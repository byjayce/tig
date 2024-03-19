package object

import "encoding/json"

type MKTagParam struct {
	TigDir     string
	ObjectType Type
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

// MKTag 태그 객체를 만드는 함수
func MKTag(param MKTagParam) (string, error) {
	body := tagBody{
		Object: param.ObjectHash,
		Type:   param.ObjectType,
		Tag:    param.Name,
		Tagger: param.Tagger,
	}

	data, err := json.Marshal(body)
	if err != nil {
		return "", err
	}

	return HashObject(HashObjectParam{
		TigDir: param.TigDir,
		Type:   Tag,
		Data:   data,
	})
}
