package wallet_manager

import (
	"context"
	"runtime"
	"sync"
	"time"

	pbCommon "github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/pkg/grpc/common"

	"github.com/crypto-bundle/bc-wallet-common-hdwallet-api/internal/app"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/anypb"
)

type unitWrapper struct {
	logger     *zap.Logger
	ctx        context.Context
	cancelFunc context.CancelFunc
	Timer      *time.Timer
	ttl        time.Duration
	Unit       WalletPoolUnitService
	notifyChan chan uuid.UUID
}

func (w *unitWrapper) Run() error {
	startedWg := &sync.WaitGroup{}
	startedWg.Add(1)

	go func(wrapped *unitWrapper, workDoneWaiter *sync.WaitGroup) {
		rawUUID, funcErr := uuid.Parse(wrapped.Unit.GetWalletUUID())
		if funcErr != nil {
			wrapped.logger.Error("unable parse wallet uuid string", zap.Error(funcErr),
				zap.String(app.WalletUUIDTag, wrapped.Unit.GetWalletUUID()))
			return
		}

		wrapped.Timer = time.NewTimer(wrapped.ttl)

		workDoneWaiter.Done()

		select {
		case fired, _ := <-wrapped.Timer.C:
			loopErr := wrapped.shutdown()
			if loopErr != nil {
				wrapped.logger.Error("unable to unload wallet data by ticker", zap.Error(loopErr),
					zap.Time(app.TickerEventTriggerTimeTag, fired))
			}

			wrapped.cancelFunc()

			break

		case <-wrapped.ctx.Done():
			loopErr := wrapped.shutdown()
			if loopErr != nil {
				wrapped.logger.Error("unable to shutdown by ctx cancel", zap.Error(loopErr))
			}

			break
		}

		wrapped.notifyChan <- rawUUID

		w.logger.Info("wallet successfully unloaded",
			zap.String(app.WalletUUIDTag, rawUUID.String()))

		return
	}(w, startedWg)

	startedWg.Wait()

	w.logger.Info("wallet successfully loaded")

	return nil
}

func (w *unitWrapper) Shutdown() {
	w.cancelFunc()
}

func (w *unitWrapper) shutdown() error {
	err := w.Unit.UnloadWallet()
	if err != nil {
		return err
	}

	w.Unit = nil
	w.Timer.Stop()
	w.Timer = nil
	w.ctx = nil

	return nil
}

func newUnitWrapper(ctx context.Context, logger *zap.Logger,
	ttl time.Duration,
	unit WalletPoolUnitService,
	notifyChan chan uuid.UUID,
) *unitWrapper {
	unitCtx, cancelFunc := context.WithCancel(ctx)

	wrapper := &unitWrapper{
		ctx:        unitCtx,
		logger:     logger,
		cancelFunc: cancelFunc,
		Timer:      nil, // will be filled in go-routine
		ttl:        ttl,
		Unit:       unit,
		notifyChan: notifyChan,
	}

	return wrapper
}

type Pool struct {
	mu     sync.Mutex
	logger *zap.Logger
	cfg    configService

	runTimeCtx context.Context

	encryptSvc      encryptService
	walletMakerFunc walletMakerFunc

	walletUnits map[uuid.UUID]*unitWrapper
	notifyChan  chan uuid.UUID
}

func (p *Pool) Run() {
	go func() {
		for {
			select {
			case <-p.runTimeCtx.Done():
				if len(p.walletUnits) != 0 {
					continue
				}

				return
			case walletUUID := <-p.notifyChan:
				p.mu.Lock()

				p.walletUnits[walletUUID] = nil
				delete(p.walletUnits, walletUUID)

				runtime.GC()

				p.mu.Unlock()
			}
		}
	}()
}

func (p *Pool) AddAndStartWalletUnit(_ context.Context,
	walletUUID uuid.UUID,
	timeToLive time.Duration,
	mnemonicEncryptedData []byte,
) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	wuWrapper, isExists := p.walletUnits[walletUUID]
	if isExists {
		wuWrapper.Timer.Reset(timeToLive)

		return nil
	}

	decryptedData, err := p.encryptSvc.Decrypt(mnemonicEncryptedData)
	if err != nil {
		return err
	}

	walletUnitInt, err := p.walletMakerFunc(walletUUID.String(), string(decryptedData))
	if err != nil {
		return err
	}

	walletUnit, isCasted := walletUnitInt.(WalletPoolUnitService)
	if !isCasted {
		return ErrUnableCastPluginEntryToPoolUnitWorker
	}

	wrapper := newUnitWrapper(p.runTimeCtx, p.logger, timeToLive, walletUnit, p.notifyChan)

	p.walletUnits[walletUUID] = wrapper

	err = wrapper.Run()
	if err != nil {
		return err
	}

	return nil
}

func (p *Pool) UnloadWalletUnit(ctx context.Context,
	mnemonicWalletUUID uuid.UUID,
) (*uuid.UUID, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	wUint, isExists := p.walletUnits[mnemonicWalletUUID]
	if !isExists {
		return nil, nil
	}
	walletUUID := wUint.Unit.GetWalletUUID()

	wUint.Shutdown()

	rawUUID, err := uuid.Parse(walletUUID)
	if err != nil {
		return nil, err
	}

	return &rawUUID, nil
}

func (p *Pool) UnloadMultipleWalletUnit(ctx context.Context,
	mnemonicWalletUUIDs []uuid.UUID,
) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	for _, v := range mnemonicWalletUUIDs {
		wUint, isExists := p.walletUnits[v]
		if !isExists {
			continue
		}

		wUint.Shutdown()
	}

	return nil
}

func (p *Pool) GetAccountAddress(ctx context.Context,
	mnemonicWalletUUID uuid.UUID,
	accountParameters *anypb.Any,
) (*string, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	wUnit, isExists := p.walletUnits[mnemonicWalletUUID]
	if !isExists {
		return nil, nil
	}

	return wUnit.Unit.GetAccountAddress(ctx, accountParameters)
}

func (p *Pool) GetMultipleAccounts(ctx context.Context,
	mnemonicWalletUUID uuid.UUID,
	multipleAccountsParameters *anypb.Any,
) (uint, []*pbCommon.AccountIdentity, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	wUnit, isExists := p.walletUnits[mnemonicWalletUUID]
	if !isExists {
		return 0, nil, nil
	}

	defer runtime.GC()

	return wUnit.Unit.GetMultipleAccounts(ctx, multipleAccountsParameters)
}

func (p *Pool) LoadAccount(ctx context.Context,
	mnemonicWalletUUID uuid.UUID,
	accountParameters *anypb.Any,
) (*string, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	wUnit, isExists := p.walletUnits[mnemonicWalletUUID]
	if !isExists {
		return nil, nil
	}

	return wUnit.Unit.LoadAccount(ctx, accountParameters)
}

func (p *Pool) SignData(ctx context.Context,
	mnemonicUUID uuid.UUID,
	accountParameters *anypb.Any,
	dataForSign []byte,
) (*string, []byte, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	wUnit, isExists := p.walletUnits[mnemonicUUID]
	if !isExists {
		p.logger.Error("wallet is not exists in wallet pool",
			zap.String(app.WalletUUIDTag, mnemonicUUID.String()))

		return nil, nil, ErrPassedWalletNotFound
	}

	return wUnit.Unit.SignData(ctx, accountParameters, dataForSign)
}

func NewWalletPool(ctx context.Context,
	logger *zap.Logger,
	cfg configService,
	mnemoWalletMakerFunc walletMakerFunc,
	encryptSrv encryptService,
) *Pool {
	return &Pool{
		runTimeCtx:      ctx,
		logger:          logger,
		cfg:             cfg,
		encryptSvc:      encryptSrv,
		walletMakerFunc: mnemoWalletMakerFunc,
		walletUnits:     make(map[uuid.UUID]*unitWrapper),
		notifyChan:      make(chan uuid.UUID),
	}
}
