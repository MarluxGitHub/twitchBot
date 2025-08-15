package service

type CommandService interface {
	HandleCommand(command string)
}

type CommandServiceImpl struct {
}

func (t *CommandServiceImpl) HandleCommand(command string) {
}

func NewCommandService() CommandService {
	return &CommandServiceImpl{}
}
