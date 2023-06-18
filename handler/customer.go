package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/clebersonp/customer-go-rest-api/model"
)

type customerHandlers struct {
	sync.Mutex
	data map[string]*model.Customer
}

func NewCustomerHandler() *customerHandlers {
	h := &customerHandlers{
		data: make(map[string]*model.Customer),
	}
	customerRoutes(h)
	return h
}

func customerRoutes(h *customerHandlers) {
	addRoutes([]route{
		newRouter("GET", "/customers/?", h.getAll),
		newRouter("GET", "/customers/([^/]+)/?", h.getCustomerById),
		newRouter("POST", "/customers/?", h.createCustomer),
		newRouter("GET", "/customers/([^/]+)/orders/?", h.getCustomerOrders),
		newRouter("POST", "/customers/([^/]+)/orders/?", h.createCustomerOrder),
		newRouter("GET", "/customers/([^/]+)/orders/([^/]+)/?", h.getCustomerOrderById),
	})
}

func (h *customerHandlers) Handlers(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s", r.Method, r.URL)
	var allow []string
	for _, route := range routes {
		matches := route.regex.FindStringSubmatch(r.URL.Path)
		log.Printf("Customer Route: {%v, %v}, Matches: %v", route.method, route.regex.String(), matches)
		if len(matches) > 0 {
			if r.Method != route.method {
				allow = append(allow, route.method)
				continue
			}
			ctx := context.WithValue(r.Context(), ctxKey{}, matches[1:])
			route.handler(w, r.WithContext(ctx))
			return
		}
	}
	if len(allow) > 0 {
		w.Header().Add("Allow", strings.Join(allow, ", "))
		writeDefaultStr(w, http.StatusMethodNotAllowed, "405 method not allowed")
		return
	}
	http.NotFound(w, r)
}

func (h *customerHandlers) getAll(w http.ResponseWriter, r *http.Request) {
	h.Lock()
	customers := make([]model.Customer, 0, len(h.data))

	for _, c := range h.data {
		customers = append(customers, *c)
	}
	h.Unlock()

	cBytes, err := json.Marshal(customers)
	if err != nil {
		writeDefaultStr(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeDefaultBytes(w, http.StatusOK, cBytes)
}

func (h *customerHandlers) getCustomerById(w http.ResponseWriter, r *http.Request) {
	id := getField(r, 0)
	h.Lock()
	c, ok := h.data[id]
	h.Unlock()
	if !ok {
		writeDefaultStr(w, http.StatusNotFound, "customer not found")
		return
	}

	cBytes, err := json.Marshal(c)
	if err != nil {
		writeDefaultStr(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeDefaultBytes(w, http.StatusOK, cBytes)
}

func (h *customerHandlers) createCustomer(w http.ResponseWriter, r *http.Request) {
	if ct := r.Header.Get("content-type"); ct != "application/json" {
		writeDefaultStr(w, http.StatusUnsupportedMediaType, fmt.Sprintf("need content-type 'application/json', but was '%s'", ct))
		return
	}

	bodyBytes, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	log.Printf("Body content %v", string(bodyBytes))
	if err != nil {
		writeDefaultStr(w, http.StatusInternalServerError, err.Error())
		return
	}
	var customerReq model.Customer
	if err := json.Unmarshal(bodyBytes, &customerReq); err != nil {
		writeDefaultStr(w, http.StatusBadRequest, err.Error())
		return
	}

	h.Lock()
	newCustomer := model.NewCustomer(customerReq.Name, customerReq.Age)
	h.data[newCustomer.Id] = &newCustomer
	h.Unlock()

	customerBytes, err := json.Marshal(newCustomer)
	if err != nil {
		writeDefaultStr(w, http.StatusInternalServerError, err.Error())
		return
	}

	formatLocation := "%s/%s"
	if strings.HasSuffix(r.URL.Path, "/") {
		formatLocation = "%s%s"
	}
	w.Header().Add("Location", fmt.Sprintf(formatLocation, r.URL.Path, newCustomer.Id))
	writeDefaultBytes(w, http.StatusCreated, customerBytes)
}

func (h *customerHandlers) getCustomerOrders(w http.ResponseWriter, r *http.Request) {
	customerId := getField(r, 0)
	h.Lock()
	c, ok := h.data[customerId]
	h.Unlock()
	if !ok {
		writeDefaultStr(w, http.StatusConflict, "customer not found")
		return
	}

	oBytes, err := json.Marshal(c.Orders)
	if err != nil {
		writeDefaultStr(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeDefaultBytes(w, http.StatusOK, oBytes)
}

func (h *customerHandlers) createCustomerOrder(w http.ResponseWriter, r *http.Request) {
	if ct := r.Header.Get("content-type"); ct != "application/json" {
		writeDefaultStr(w, http.StatusUnsupportedMediaType, fmt.Sprintf("need content-type 'application/json', but was '%s'", ct))
		return
	}

	customerId := getField(r, 0)
	customer, ok := h.data[customerId]
	if !ok {
		writeDefaultStr(w, http.StatusConflict, "customer not found")
		return
	}

	bodyBytes, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	log.Printf("Body content %v", string(bodyBytes))
	if err != nil {
		writeDefaultStr(w, http.StatusInternalServerError, err.Error())
		return
	}
	var orderReq model.Order
	if err := json.Unmarshal(bodyBytes, &orderReq); err != nil {
		writeDefaultStr(w, http.StatusBadRequest, err.Error())
		return
	}

	h.Lock()
	newOrder := model.NewOrder(orderReq.Product, orderReq.Price, orderReq.Amount)
	customer.Add(newOrder)
	h.Unlock()

	orderBytes, err := json.Marshal(newOrder)
	if err != nil {
		writeDefaultStr(w, http.StatusInternalServerError, err.Error())
		return
	}

	formatLocation := "%s/%s"
	if strings.HasSuffix(r.URL.Path, "/") {
		formatLocation = "%s%s"
	}
	w.Header().Add("Location", fmt.Sprintf(formatLocation, r.URL.Path, newOrder.Id))
	writeDefaultBytes(w, http.StatusCreated, orderBytes)
}

func (h *customerHandlers) getCustomerOrderById(w http.ResponseWriter, r *http.Request) {
	customerId := getField(r, 0)
	h.Lock()
	customer, ok := h.data[customerId]
	h.Unlock()
	if !ok {
		writeDefaultStr(w, http.StatusConflict, "customer not found")
		return
	}

	orderId := getField(r, 1)
	order, exists := customer.Contains(orderId)
	if !exists {
		writeDefaultStr(w, http.StatusNotFound, "order not found")
		return
	}

	oBytes, err := json.Marshal(order)
	if err != nil {
		writeDefaultStr(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeDefaultBytes(w, http.StatusOK, oBytes)
}

func getField(r *http.Request, index int) string {
	fields := r.Context().Value(ctxKey{}).([]string)
	if index >= 0 && index < len(fields) {
		return fields[index]
	}
	return ""
}
