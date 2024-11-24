package keto

import (
	"context"
	"fmt"

	"github.com/dupmanio/dupman/packages/common/otel"
	ory "github.com/ory/client-go"
	"go.uber.org/zap"
)

type Keto struct {
	logger *zap.Logger
	ot     *otel.OTel
	config Config

	writeClient *ory.APIClient
}

func New(conf Config, logger *zap.Logger, ot *otel.OTel) (*Keto, error) {
	writeConfig := ory.NewConfiguration()
	writeConfig.Servers = []ory.ServerConfiguration{
		{
			URL: conf.WriteURL,
		},
	}

	return &Keto{
		logger:      logger,
		ot:          ot,
		config:      conf,
		writeClient: ory.NewAPIClient(writeConfig),
	}, nil
}

func (inst *Keto) CreateRelationship(
	ctx context.Context,
	namespace string,
	object string,
	relation string,
	subjectID string,
) error {
	ctx, span := inst.ot.GetSpanForFunctionCall(
		ctx,
		1,
		otel.OryKetoNamespace(namespace),
		otel.OryKetoObject(object),
		otel.OryKetoRelation(relation),
		otel.OryKetoSubjectID(subjectID),
	)
	defer span.End()

	payload := ory.CreateRelationshipBody{
		Namespace: &namespace,
		Object:    &object,
		Relation:  &relation,
		SubjectId: &subjectID,
	}

	_, _, err := inst.
		writeClient.
		RelationshipAPI.
		CreateRelationship(ctx).
		CreateRelationshipBody(payload).
		Execute()
	if err != nil {
		return fmt.Errorf("unable to create relationship: %w", err)
	}

	return nil
}

func (inst *Keto) DeleteRelationship(
	ctx context.Context,
	namespace string,
	object string,
	relation string,
	subjectID string,
) error {
	ctx, span := inst.ot.GetSpanForFunctionCall(
		ctx,
		1,
		otel.OryKetoNamespace(namespace),
		otel.OryKetoObject(object),
		otel.OryKetoRelation(relation),
		otel.OryKetoSubjectID(subjectID),
	)
	defer span.End()

	_, err := inst.
		writeClient.
		RelationshipAPI.
		DeleteRelationships(ctx).
		Namespace(namespace).
		Object(object).
		Relation(relation).
		SubjectId(subjectID).
		Execute()
	if err != nil {
		return fmt.Errorf("unable to delete relationship: %w", err)
	}

	return nil
}
