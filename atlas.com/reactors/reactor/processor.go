package reactor

import (
	"atlas-reactors/data/reactor"
	"atlas-reactors/kafka/producer"
	"context"
	"github.com/Chronicle20/atlas-model/model"
	"github.com/Chronicle20/atlas-tenant"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
)

type Processor struct {
	l   logrus.FieldLogger
	ctx context.Context
	t   tenant.Model
	rd  *reactor.Processor
}

func NewProcessor(l logrus.FieldLogger, ctx context.Context) *Processor {
	p := &Processor{
		l:   l,
		ctx: ctx,
		t:   tenant.MustFromContext(ctx),
		rd:  reactor.NewProcessor(l, ctx),
	}
	return p
}

func (p *Processor) GetById(id uint32) (Model, error) {
	return GetRegistry().Get(id)
}

func (p *Processor) GetInMap(worldId byte, channelId byte, mapId uint32) ([]Model, error) {
	return GetRegistry().GetInMap(p.t, worldId, channelId, mapId), nil
}

func (p *Processor) Create(b *ModelBuilder) error {
	d, err := p.rd.GetById(b.Classification())
	if err != nil {
		p.l.WithError(err).Errorf("Unable to retrieve reactor [%d] game data.", b.Classification())
		return err
	}
	b.SetData(d)
	r := GetRegistry().Create(p.t, b)
	p.l.Debugf("Created reactor [%d] of [%d].", r.Id(), r.Classification())
	return producer.ProviderImpl(p.l)(p.ctx)(EnvEventStatusTopic)(createdStatusEventProvider(r))
}

func Teardown(l logrus.FieldLogger) func() {
	return func() {
		sctx, span := otel.GetTracerProvider().Tracer("atlas-reactors").Start(context.Background(), "teardown")
		defer span.End()

		err := model.ForEachMap(model.FixedProvider(GetRegistry().GetAll()), func(t tenant.Model) model.Operator[[]Model] {
			tctx := tenant.WithContext(sctx, t)
			return model.ExecuteForEachSlice(NewProcessor(l, tctx).Destroy, model.ParallelExecute())
		}, model.ParallelExecute())
		if err != nil {
			l.WithError(err).Errorf("Error destroying all reactors on teardown.")
		}
	}
}

func (p *Processor) Destroy(m Model) error {
	GetRegistry().Remove(p.t, m.Id())
	return producer.ProviderImpl(p.l)(p.ctx)(EnvEventStatusTopic)(destroyedStatusEventProvider(m))
}
