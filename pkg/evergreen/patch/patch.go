package patch

type Body struct {
	Priority int `json:"priority"`
}

type Patch struct {
	PatchId     string `json:"patch_id"`
	ProjectId   string `json:"project_id"`
	Description string `json:"description"`
}
