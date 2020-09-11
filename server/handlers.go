package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"mime"
	"net/http"
	"strconv"

	"github.com/rs/zerolog"
	"github.com/valyala/fastjson"
	hl "github.com/vishal1132/cafebucks/handlers"
)

type beans struct {
	Beans   string `json:"beans"`
	Stock   int    `json:"stock,omitempty"`
	UnitUse int    `json:"unituse,omitempty"`
}

type handler struct {
	l *zerolog.Logger
}

type ctxKey uint8

const maxBodySize = 2 * 1024 * 1024 // 2MB

const (
	ctxKeyReqID ctxKey = iota
)

func (s *server) registerHandlers() {
	h := handler{l: &s.logger}
	s.mux.HandleFunc("/beansService/_health/status", h.handleHealth).Methods("GET")
	s.mux.HandleFunc("/stock", h.handleGetStock).Methods("GET")
	s.mux.HandleFunc("/addBeans", h.handleAddBeans).Methods("POST")
}

func (h *handler) handleHealth(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "healthy")
}

func validateBeans(coffee string) bool {
	if val, ok := coffeebeans[coffee]; ok {
		if val.Stock > val.UnitUse {
			val.Stock = val.Stock - val.UnitUse
			return true
		}
		return false
	}
	return false
}

func (h *handler) handleGetStock(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	l := h.l.With().Str("context", "get stock event")
	rid, ok := hl.CtxRequestID(ctx)
	if ok {
		l = l.Str("request_id", rid)
	}

	lg := l.Logger()

	if len(beanSlice) == 0 {
		lg.Info().
			Msg("No beans in the stock")
		io.WriteString(w, "No beans in our cafę")
		return
	}

	for id, val := range beanSlice {
		lg.Info().
			Str("id", strconv.Itoa(id)).
			Str("beans", val.Beans).
			Str("stock", strconv.Itoa(val.Stock)).
			Str("unit", strconv.Itoa(val.UnitUse)).
			Msg("Beans")
		io.WriteString(w, fmt.Sprintf("beans: %s, stock: %d, Unit Consumption: %d\n", val.Beans, val.Stock, val.UnitUse))
	}
}

func (h *handler) handleAddBeans(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	l := h.l.With().Str("context", "add beans event")
	rid, ok := hl.CtxRequestID(ctx)
	if ok {
		l = l.Str("request_id", rid)
	}

	lg := l.Logger()
	mtype, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil {
		lg.Error().
			Err(err).
			Msg("Failed to parse content type")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if mtype != "application/json" {
		lg.Error().
			Str("content_type", mtype).
			Msg("content type was not JSON")
		io.WriteString(w, "Not Json Type")
		w.Header().Set("Accept", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// Now unmarshaling the body into Coffee struct from cafebucks/eventbus
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, maxBodySize))
	if err != nil {
		lg.Error().
			Err(err).
			Msg("failed to read request body")

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	document, err := fastjson.ParseBytes(body)
	if err != nil {
		lg.Error().
			Err(err).
			Msg("failed to unmarshal JSON document")

		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	beanReq, err := hl.GetJSONString(document, "beans")
	if err != nil {
		lg.Error().
			Err(err).
			Msg("no beans field in request body")
		io.WriteString(w, "No beans field in the request body")
		return
	}
	quantity, err := hl.GetJSONInt(document, "quantity")
	if err != nil {
		lg.Error().
			Err(err).
			Msg("quantity field wrong in request body")
		io.WriteString(w, "quantity field wrong in the request body")
		return
	}

	id, stock, unit, exists := checkExist(beanReq)

	if exists {
		beanSlice[id].Stock += quantity
		lg.Info().
			Str("stock", strconv.Itoa(stock+quantity)).
			Str("unit", strconv.Itoa(unit)).
			Msg("Stock Updated")
		io.WriteString(w, "Stock updated")
		return
	}
	// beans does not exist, create a new entry
	unitQty, err := hl.GetJSONInt(document, "unit")
	if err != nil {
		lg.Error().
			Err(err).
			Msg("failed to update beans no unit field in request body")
		io.WriteString(w, "No unit field in request body")
	}
	var beanVar = beans{Beans: beanReq, Stock: quantity, UnitUse: unitQty}

	beanSlice = append(beanSlice, &beanVar)
	lg.Info().
		Str("beans", beanReq).
		Str("stock", strconv.Itoa(beanVar.Stock)).
		Str("unit", strconv.Itoa(beanVar.UnitUse)).
		Msg("Stock Updated, added new beans")
	io.WriteString(w, "stock updated, added new beans")
	return
}

func checkExist(beans string) (int, int, int, bool) {
	for id, val := range beanSlice {
		if val.Beans == beans {
			return id, val.Stock, val.UnitUse, true
		}
	}
	return -1, 0, 0, false
}
