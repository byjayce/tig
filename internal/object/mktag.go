package object

import "encoding/json"

type MKTagParam struct {
	BaseDir    string
	ObjectHash string
	Name       string
	Message    string
	Tagger     string
}

type tagBody struct {
	Object string
	Type   Type
	Tag    string
	Tagger string
}

func MKTag(param MKTagParam) (string, error) {
	body := tagBody{
		Object: param.ObjectHash,
		Type:   TagType,
		Tag:    param.Name,
		Tagger: param.Tagger,
	}

	data, err := json.Marshal(body)
	if err != nil {
		return "", err
	}

	return HashObject(HashObjectParam{
		BaseDir: param.BaseDir,
		Type:    TagType,
		Data:    data,
	})
}
