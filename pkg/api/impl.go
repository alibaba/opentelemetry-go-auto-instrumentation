// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package api

// For testing purpose only
type CallContextImpl struct {
	Params     []interface{}
	ReturnVals []interface{}
	SkipCall   bool
	Data       interface{}
}

func (c *CallContextImpl) SetSkipCall(skip bool)    { c.SkipCall = skip }
func (c *CallContextImpl) IsSkipCall() bool         { return c.SkipCall }
func (c *CallContextImpl) SetData(data interface{}) { c.Data = data }
func (c *CallContextImpl) GetData() interface{}     { return c.Data }
func (c *CallContextImpl) GetKeyData(key string) interface{} {
	if c.Data == nil {
		return nil
	}
	return c.Data.(map[string]interface{})[key]
}
func (c *CallContextImpl) SetKeyData(key string, val interface{}) {
	if c.Data == nil {
		c.Data = make(map[string]interface{})
	}
	c.Data.(map[string]interface{})[key] = val
}

func (c *CallContextImpl) HasKeyData(key string) bool {
	if c.Data == nil {
		return false
	}
	_, ok := c.Data.(map[string]interface{})[key]
	return ok
}

func (c *CallContextImpl) GetParam(idx int) interface{} {
	return &c.Params[idx]
}

func (c *CallContextImpl) SetParam(idx int, val interface{}) {
	c.Params[idx] = val
}

func (c *CallContextImpl) GetReturnVal(idx int) interface{} {
	return &c.ReturnVals[idx]
}

func (c *CallContextImpl) SetReturnVal(idx int, val interface{}) {
	c.ReturnVals[idx] = val
}

func NewCallContext() CallContext {
	return &CallContextImpl{
		Params:     make([]interface{}, 1024),
		ReturnVals: make([]interface{}, 1024),
		SkipCall:   false,
	}
}
