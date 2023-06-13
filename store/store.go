package store

import (
	"context"
	"encoding/json"
	"errors"
	tracer "example.com/mod/tracer"
	"fmt"
	"github.com/google/uuid"
	"github.com/hashicorp/consul/api"
	"os"
)

type Store struct {
	cli *api.Client
}

func New() (*Store, error) {
	db := os.Getenv("DB")
	dbport := os.Getenv("DBPORT")

	config := api.DefaultConfig()
	config.Address = fmt.Sprintf("%s:%s", db, dbport)
	client, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}

	return &Store{
		cli: client,
	}, nil
}

func (ps *Store) Get(ctx context.Context, id string, version string) ([]*Config, error) {
	span := tracer.StartSpanFromContext(ctx, "GetConfig")
	defer span.Finish()
	kv := ps.cli.KV()

	data, _, err := kv.List(constructKey(id, version, ""), nil)
	if err != nil {
		tracer.LogError(span, err)
		return nil, err
	}

	configs := []*Config{}
	for _, pair := range data {
		config := &Config{}
		err = json.Unmarshal(pair.Value, config)
		if err != nil {
			tracer.LogError(span, err)
			return nil, err
		}
		configs = append(configs, config)
	}
	return configs, nil
}
func (ps *Store) GetGroup(ctx context.Context, id string, version string) ([]*Group, error) {
	span := tracer.StartSpanFromContext(ctx, "GetGroup")
	defer span.Finish()

	kv := ps.cli.KV()

	data, _, err := kv.List(constructGroupKey(id, version), nil)
	if err != nil {
		tracer.LogError(span, err)
		return nil, err
	}

	configs := []*Group{}
	for _, pair := range data {
		config := &Group{}
		err = json.Unmarshal(pair.Value, config)
		if err != nil {
			tracer.LogError(span, err)
			return nil, err
		}
		configs = append(configs, config)
	}
	return configs, nil
}
func (ps *Store) GetGroupId(ctx context.Context, id string) ([]*Group, error) {
	span := tracer.StartSpanFromContext(ctx, "GetGroup")
	defer span.Finish()

	kv := ps.cli.KV()

	data, _, err := kv.List(constructGroupKey2(id), nil)
	if err != nil {
		tracer.LogError(span, err)
		return nil, err
	}

	configs := []*Group{}
	for _, pair := range data {
		config := &Group{}
		err = json.Unmarshal(pair.Value, config)
		if err != nil {
			tracer.LogError(span, err)
			return nil, err
		}
		configs = append(configs, config)
	}
	return configs, nil
}
func (ps *Store) GetOneGroup(ctx context.Context, id string, version string) (*Group, error) {
	span := tracer.StartSpanFromContext(ctx, "GetOneGroup")
	defer span.Finish()
	kv := ps.cli.KV()

	data, _, err := kv.List(constructGroupKey(id, version), nil)
	if err != nil {
		tracer.LogError(span, err)
		return nil, err
	}

	for _, pair := range data {
		config := &Group{}
		err = json.Unmarshal(pair.Value, config)
		if err != nil {
			tracer.LogError(span, err)
			return nil, err
		}
		return config, nil
	}

	// Ako nijedna grupa nije pronađena, možete vratiti odgovarajuću grešku
	return nil, errors.New("group not found, not exist")
}
func (ps *Store) GetOneGroup2(ctx context.Context, id string) (*Group, error) {
	span := tracer.StartSpanFromContext(ctx, "GetOneGroup")
	defer span.Finish()
	kv := ps.cli.KV()

	data, _, err := kv.List(constructGroupKey2(id), nil)
	if err != nil {
		tracer.LogError(span, err)
		return nil, err
	}

	for _, pair := range data {
		config := &Group{}
		err = json.Unmarshal(pair.Value, config)
		if err != nil {
			tracer.LogError(span, err)
			return nil, err
		}
		return config, nil
	}

	// Ako nijedna grupa nije pronađena, možete vratiti odgovarajuću grešku
	return nil, errors.New("group not found, not exist")
}
func (ps *Store) GetOneConfig(ctx context.Context, id string, version string) (*Config, error) {
	span := tracer.StartSpanFromContext(ctx, "GetOneConfig")
	defer span.Finish()
	kv := ps.cli.KV()

	data, _, err := kv.List(constructKey(id, version, ""), nil)
	if err != nil {
		tracer.LogError(span, err)
		return nil, err
	}

	for _, pair := range data {
		config := &Config{}
		err = json.Unmarshal(pair.Value, config)
		if err != nil {
			tracer.LogError(span, err)
			return nil, err
		}
		return config, nil
	}

	// Ako nijedna grupa nije pronađena, možete vratiti odgovarajuću grešku
	return nil, errors.New("config not found, not exist")
}
func (ps *Store) GetOneConfig2(ctx context.Context, id string) (*Config, error) {
	span := tracer.StartSpanFromContext(ctx, "GetOneConfig")
	defer span.Finish()
	kv := ps.cli.KV()

	data, _, err := kv.List(constructKey2(id), nil)
	if err != nil {
		tracer.LogError(span, err)
		return nil, err
	}

	for _, pair := range data {
		config := &Config{}
		err = json.Unmarshal(pair.Value, config)
		if err != nil {
			tracer.LogError(span, err)
			return nil, err
		}
		return config, nil
	}

	// Ako nijedna grupa nije pronađena, možete vratiti odgovarajuću grešku
	return nil, errors.New("config not found, not exist")
}

func (ps *Store) SaveGroup(ctx context.Context, post *Group) (*Group, error) {
	span := tracer.StartSpanFromContext(ctx, "SaveGroup")
	defer span.Finish()
	kv := ps.cli.KV()

	/*sid, rid := generateGroupKey(post.Version, post.Labels)
	post.Id = rid
	*/
	data, err := json.Marshal(post)
	if err != nil {
		tracer.LogError(span, err)
		return nil, err
	}

	p := &api.KVPair{Key: constructGroupKey(post.Id, post.Version), Value: data}
	_, err = kv.Put(p, nil)
	if err != nil {
		tracer.LogError(span, err)
		return nil, err
	}

	return post, nil
}

func (ps *Store) GetAll(ctx context.Context) ([]*Config, error) {
	span := tracer.StartSpanFromContext(ctx, "GetAll")
	defer span.Finish()
	kv := ps.cli.KV()
	data, _, err := kv.List(all, nil)
	if err != nil {
		tracer.LogError(span, err)
		return nil, err
	}

	configs := []*Config{}
	for _, pair := range data {
		config := &Config{}
		err = json.Unmarshal(pair.Value, config)
		if err != nil {
			tracer.LogError(span, err)
			return nil, err
		}
		configs = append(configs, config)
	}

	return configs, nil
}
func (ps *Store) GetAllGroups(ctx context.Context) ([]*Group, error) {
	span := tracer.StartSpanFromContext(ctx, "GetAllGroups")
	defer span.Finish()
	kv := ps.cli.KV()
	data, _, err := kv.List(allGroups, nil)
	if err != nil {
		tracer.LogError(span, err)
		return nil, err
	}

	groups := []*Group{}
	for _, pair := range data {
		group := &Group{}
		err = json.Unmarshal(pair.Value, group)
		if err != nil {
			tracer.LogError(span, err)
			return nil, err
		}
		groups = append(groups, group)
	}

	return groups, nil
}

func (ps *Store) Delete(ctx context.Context, id string, version string) (map[string]string, error) {
	span := tracer.StartSpanFromContext(ctx, "Delete")
	defer span.Finish()
	kv := ps.cli.KV()
	_, err := kv.DeleteTree(constructKey(id, version, ""), nil)
	if err != nil {
		tracer.LogError(span, err)
		return nil, err
	}

	return map[string]string{"Deleted": id}, nil
}
func (ps *Store) DeleteByLabel(ctx context.Context, id string, version string, labels string) (map[string]string, error) {
	span := tracer.StartSpanFromContext(ctx, "Delete")
	defer span.Finish()
	kv := ps.cli.KV()
	_, err := kv.DeleteTree(constructKey(id, version, labels), nil)
	if err != nil {
		tracer.LogError(span, err)
		return nil, err
	}

	return map[string]string{"Deleted": id}, nil
}
func (ps *Store) DeleteGroup(ctx context.Context, id string, version string) (map[string]string, error) {
	span := tracer.StartSpanFromContext(ctx, "DeleteGroup")
	defer span.Finish()
	kv := ps.cli.KV()
	_, err := kv.DeleteTree(constructGroupKey(id, version), nil)
	if err != nil {
		tracer.LogError(span, err)
		return nil, err
	}

	return map[string]string{"Deleted": id}, nil
}

func (ps *Store) DeleteGroupId(ctx context.Context, id string) (map[string]string, error) {
	span := tracer.StartSpanFromContext(ctx, "DeleteGroup")
	defer span.Finish()
	kv := ps.cli.KV()
	_, err := kv.DeleteTree(constructGroupKey2(id), nil)
	if err != nil {
		tracer.LogError(span, err)
		return nil, err
	}

	return map[string]string{"Deleted": id}, nil
}

func (ps *Store) Config(ctx context.Context, config *Config) (*Config, error) {
	span := tracer.StartSpanFromContext(ctx, "Config")
	defer span.Finish()

	kv := ps.cli.KV()

	sid, rid := generateKey(config.Version, config.Labels)
	config.Id = rid

	data, err := json.Marshal(config)
	if err != nil {
		tracer.LogError(span, err)
		return nil, err
	}

	p := &api.KVPair{Key: sid, Value: data}
	_, err = kv.Put(p, nil)
	if err != nil {
		tracer.LogError(span, err)
		return nil, err
	}

	return config, nil
}

func (ps *Store) PostGroup(ctx context.Context, post *Group) (*Group, error) {
	span := tracer.StartSpanFromContext(ctx, "PostGroup")
	defer span.Finish()

	kv := ps.cli.KV()

	sid, rid := generateGroupKey(post.Version)
	post.Id = rid

	data, err := json.Marshal(post)
	if err != nil {
		tracer.LogError(span, err)
		return nil, err
	}

	p := &api.KVPair{Key: sid, Value: data}
	_, err = kv.Put(p, nil)
	if err != nil {
		tracer.LogError(span, err)
		return nil, err
	}

	return post, nil
}

/*
	func (ps *Store) GetGroupsByLabels(ctx context.Context, id string, version string, labels string) ([]*Group, error) {
		span := tracer.StartSpanFromContext(ctx, "GetGroupsByLabel")
		defer span.Finish()
		kv := ps.cli.KV()

		data, _, err := kv.List(constructGroupKey(id, version,labels), nil)
		if err != nil {
			tracer.LogError(span, err)
			return nil, err
		}

		posts := []*Group{}

		for _, pair := range data {
			post := &Group{}
			err = json.Unmarshal(pair.Value, post)
			if err != nil {
				tracer.LogError(span, err)
				return nil, err
			}
			posts = append(posts, post)
		}

		if err != nil {
			tracer.LogError(span, err)
			return nil, err
		}

		return posts, nil
	}
*/
func (ps *Store) GetConfigsByLabels(ctx context.Context, id string, version string, labels string) ([]*Config, error) {
	span := tracer.StartSpanFromContext(ctx, "GetConfigsByLabel")
	defer span.Finish()
	kv := ps.cli.KV()

	data, _, err := kv.List(constructKey(id, version, labels), nil)
	if err != nil {
		tracer.LogError(span, err)
		return nil, err
	}

	configs := []*Config{}

	for _, pair := range data {
		config := &Config{}
		err = json.Unmarshal(pair.Value, config)
		if err != nil {
			tracer.LogError(span, err)
			return nil, err
		}
		configs = append(configs, config)
	}

	if err != nil {
		tracer.LogError(span, err)
		return nil, err
	}

	return configs, nil
}

func (ps *Store) SaveRequestId(ctx context.Context) string {
	span := tracer.StartSpanFromContext(ctx, "SaveRequestId")
	defer span.Finish()
	kv := ps.cli.KV()

	reqId := generateRequestId(ctx)

	i := &api.KVPair{Key: reqId, Value: nil}

	_, err := kv.Put(i, nil)
	if err != nil {
		tracer.LogError(span, err)
		return "error"
	}
	return reqId

}

func (ps *Store) FindRequestId(ctx context.Context, requestId string) bool {
	span := tracer.StartSpanFromContext(ctx, "GetRequestId")
	defer span.Finish()

	kv := ps.cli.KV()

	key, _, err := kv.Get(requestId, nil)

	fmt.Println(key)

	if err != nil || key == nil {
		tracer.LogError(span, err)
		return false
	}

	return true

}
func generateRequestId(ctx context.Context) string {
	span := tracer.StartSpanFromContext(ctx, "Generate")
	defer span.Finish()
	rid := uuid.New().String()

	return rid
}
