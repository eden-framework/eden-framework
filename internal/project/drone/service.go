package drone

type PipelineService struct {
	Name  string `yaml:"name" json:"name"`
	Image string `yaml:"image" json:"image"`
}

func NewPipelineService() *PipelineService {
	return new(PipelineService)
}

func (s *PipelineService) WithName(n string) *PipelineService {
	s.Name = n
	return s
}

func (s *PipelineService) WithImage(img string) *PipelineService {
	s.Image = img
	return s
}
