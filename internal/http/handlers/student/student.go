package student

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/Shreya20002/students-go/internal/storage"
	"github.com/Shreya20002/students-go/internal/types"
	"github.com/Shreya20002/students-go/internal/utils/response"
	"github.com/go-playground/validator/v10"
)

func New(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("creating a student")
		// json info received , so we will create a struct of that format to receive such data in go
		var student types.Student

		err := json.NewDecoder(r.Body).Decode(&student) // decode the json data into the struct var student
		if errors.Is(err, io.EOF) {
			// 404 bad request waala error StatusBadRequest
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("empty body")))
			return
		}

		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		// for production lvl code, we need to validate the errors
		if err := validator.New().Struct(student); err != nil {

			validateErrs := err.(validator.ValidationErrors)
			response.WriteJson(w, http.StatusBadRequest, response.ValidationError(validateErrs))
			return
		}
		lastId, err := storage.CreateStudent(
			student.Name,
			student.Email,
			student.Age,
		)

		slog.Info("user created successfully", slog.String("userId", fmt.Sprint(lastId)))

		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, err)
			return
		}

		response.WriteJson(w, http.StatusCreated, map[string]int64{"id": lastId})
	}

	//response.WriteJson(w, http.StatusCreated, map[string]string{"success": "ok"})
	//w.Write([]byte("welcome to students api"))

}
