package handler

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/wintltr/vand-interview-crud-project/authentication"
	"github.com/wintltr/vand-interview-crud-project/model"
	"github.com/wintltr/vand-interview-crud-project/util"
)

func AddStoreHandler(w http.ResponseWriter, r *http.Request) {

	var store model.Store

	tokenData, err := authentication.ExtractTokenMetadata(r)
	if err != nil {
		util.ERROR(w, http.StatusBadRequest, errors.New("Please login").Error())
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		util.ERROR(w, http.StatusBadRequest, errors.New("Fail to read user data").Error())
		return
	}
	
	err = json.Unmarshal(body, &store)
	if err != nil {
		util.ERROR(w, http.StatusBadRequest, errors.New("Fail to parse json format").Error())
		return
	}
	store.UserId = tokenData.Userid

	err = model.InsertStoreToDB(store)

	if err != nil {
		util.JSON(w, http.StatusBadRequest, errors.New(err.Error()))
	} else {
		util.JSON(w, http.StatusOK, store)
	}
}

func RemoveStoreHandler(w http.ResponseWriter, r *http.Request) {

	tokenData, err := authentication.ExtractTokenMetadata(r)
	if err != nil {
		util.ERROR(w, http.StatusBadRequest, errors.New("Please login").Error())
		return
	}

	// Retrieve Id
	vars := mux.Vars(r)
	storeId, err := strconv.Atoi(vars["id"])
	if err != nil {
		util.ERROR(w, http.StatusUnauthorized, errors.New("Failed to retrieve Id").Error())
		return
	}

	//Check authorization
	tmpStore, err := model.GetStoreByIdFromDB(storeId)
	if tmpStore.UserId != tokenData.Userid || err != nil {
		util.ERROR(w, http.StatusBadRequest, errors.New("You are not authorized to modify this store").Error())
		return
	}

	err = model.DeleteStoreFromDB(storeId)
	if err != nil {
		util.ERROR(w, http.StatusOK, errors.New("Failed to remove store").Error())
	} else {
		util.JSON(w, http.StatusOK, nil)
	}
}

func ListStoreByPageHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	page, err := strconv.Atoi(vars["page"])
	if err != nil {
		util.ERROR(w, http.StatusBadRequest, "failed to get page value")
		return
	}

	storeList, err := model.GetStoreByPage(page, 10)

	// Return json
	if err != nil {
		util.ERROR(w, http.StatusBadRequest, errors.New("fail to get stores").Error())
	} else {
		util.JSON(w, http.StatusOK, storeList)
	}
}

func StoreDetailHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	storeId, err := strconv.Atoi(vars["id"])
	if err != nil {
		util.ERROR(w, http.StatusBadRequest, "failed to get store detail")
		return
	}

	product, err := model.GetStoreByIdFromDB(storeId)

	// Return json
	if err != nil {
		util.ERROR(w, http.StatusBadRequest, err.Error())
	} else {
		util.JSON(w, http.StatusOK, product)
	}

}

func UpdateStoreHandler(w http.ResponseWriter, r *http.Request) {

	var store model.Store

	tokenData, err := authentication.ExtractTokenMetadata(r)
	if err != nil {
		util.ERROR(w, http.StatusBadRequest, errors.New("Please login").Error())
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		util.ERROR(w, http.StatusBadRequest, errors.New("Fail to update store").Error())
		return
	}
	
	err = json.Unmarshal(body, &store)
	if err != nil {
		util.ERROR(w, http.StatusBadRequest, errors.New("Fail to parse json format").Error())
		return
	}

	//Check authorization
	tmpStore, err := model.GetStoreByIdFromDB(store.StoreId)
	if tmpStore.UserId != tokenData.Userid || err != nil {
		util.ERROR(w, http.StatusBadRequest, errors.New("You are not authorized to modify this store").Error())
		return
	}

	err = model.UpdateStoreToDB(store)

	if err != nil {
		util.JSON(w, http.StatusBadRequest, errors.New(err.Error()))
	} else {
		util.JSON(w, http.StatusOK, store)
	}
}

func SearchStoreByPageHandler(w http.ResponseWriter, r *http.Request) {

	type received struct {
		SearchString string `json:"search_string"`
		Page int `json:"page"`
	}

	var receivedValue received

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		util.ERROR(w, http.StatusBadRequest, errors.New("Fail to read user data").Error())
		return
	}
	
	err = json.Unmarshal(body, &receivedValue)
	if err != nil {
		util.ERROR(w, http.StatusBadRequest, errors.New("Fail to parse json format").Error())
		return
	}

	if receivedValue.Page == 0 {
		receivedValue.Page = 1
	}

	productList, err := model.SearchStoreByPage(receivedValue.Page, 10, receivedValue.SearchString)

	// Return json
	if err != nil {
		util.ERROR(w, http.StatusBadRequest, errors.New("fail to get products").Error())
	} else {
		util.JSON(w, http.StatusOK, productList)
	}
}
