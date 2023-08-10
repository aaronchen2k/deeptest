package service

import (
	"encoding/json"
	valueUtils "github.com/aaronchen2k/deeptest/internal/agent/exec/utils/value"
	"github.com/aaronchen2k/deeptest/internal/pkg/consts"
	"github.com/aaronchen2k/deeptest/internal/pkg/domain"
	repo "github.com/aaronchen2k/deeptest/internal/server/modules/repo"
)

type ExecConditionService struct {
	PreConditionRepo  *repo.PreConditionRepo  `inject:""`
	PostConditionRepo *repo.PostConditionRepo `inject:""`

	ExtractorRepo  *repo.ExtractorRepo  `inject:""`
	CheckpointRepo *repo.CheckpointRepo `inject:""`
	ScriptRepo     *repo.ScriptRepo     `inject:""`

	ShareVarService *ShareVarService `inject:""`
}

func (s *ExecConditionService) SavePreConditionResult(invokeId,
	debugInterfaceId, caseInterfaceId, endpointInterfaceId, serveId, processorId, scenarioId uint, usedBy consts.UsedBy,
	preConditions []domain.InterfaceExecCondition) (err error) {

	for _, condition := range preConditions {
		if condition.Type == consts.ConditionTypeScript {
			var scriptBase domain.ScriptBase
			json.Unmarshal(condition.Raw, &scriptBase)
			if scriptBase.Disabled {
				continue
			}

			scriptBase.InvokeId = invokeId

			s.ScriptRepo.UpdateResult(scriptBase)
			s.ScriptRepo.CreateLog(scriptBase)

			for _, settings := range scriptBase.VariableSettings {
				s.ShareVarService.Save(settings.Name, valueUtils.InterfaceToStr(settings.Value),
					invokeId, debugInterfaceId, caseInterfaceId, endpointInterfaceId, serveId, processorId, scenarioId,
					consts.Public, usedBy)
			}
		}
	}

	return
}

func (s *ExecConditionService) SavePostConditionResult(invokeId,
	debugInterfaceId, caseInterfaceId, endpointInterfaceId, serveId, processorId, scenarioId uint, usedBy consts.UsedBy,
	postConditions []domain.InterfaceExecCondition) (err error) {

	for _, condition := range postConditions {
		if condition.Type == consts.ConditionTypeExtractor {
			var extractorBase domain.ExtractorBase
			json.Unmarshal(condition.Raw, &extractorBase)
			if extractorBase.Disabled {
				continue
			}

			extractorBase.InvokeId = invokeId

			s.ExtractorRepo.UpdateResult(extractorBase)
			s.ExtractorRepo.CreateLog(extractorBase)

			if extractorBase.ResultStatus == consts.Pass {
				s.ShareVarService.Save(extractorBase.Variable, extractorBase.Result,
					invokeId, debugInterfaceId, caseInterfaceId, endpointInterfaceId, serveId, processorId, scenarioId,
					extractorBase.Scope, usedBy)
			}

		} else if condition.Type == consts.ConditionTypeCheckpoint {
			var checkpointBase domain.CheckpointBase
			json.Unmarshal(condition.Raw, &checkpointBase)
			if checkpointBase.Disabled {
				continue
			}

			checkpointBase.InvokeId = invokeId

			s.CheckpointRepo.UpdateResult(checkpointBase)
			s.CheckpointRepo.CreateLog(checkpointBase)

		} else if condition.Type == consts.ConditionTypeScript {
			var scriptBase domain.ScriptBase
			json.Unmarshal(condition.Raw, &scriptBase)
			if scriptBase.Disabled {
				continue
			}

			scriptBase.InvokeId = invokeId

			s.ScriptRepo.UpdateResult(scriptBase)
			s.ScriptRepo.CreateLog(scriptBase)

			for _, settings := range scriptBase.VariableSettings {
				s.ShareVarService.Save(settings.Name, valueUtils.InterfaceToStr(settings.Value),
					invokeId, debugInterfaceId, caseInterfaceId, endpointInterfaceId, serveId, processorId, scenarioId,
					consts.Public, usedBy)
			}
		}
	}

	return
}
