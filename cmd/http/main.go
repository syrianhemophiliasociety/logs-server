package main

import (
	"net/http"
	"regexp"
	"shs/actions"
	"shs/app"
	"shs/config"
	"shs/handlers/apis"
	"shs/handlers/middlewares/auth"
	"shs/handlers/middlewares/contenttype"
	"shs/handlers/middlewares/logger"
	"shs/jwt"
	"shs/log"
	"shs/mariadb"
	"shs/redis"

	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/json"
)

func main() {
	repo, err := mariadb.New()
	if err != nil {
		log.Fatalln(err)
	}
	cache := redis.New()
	app := app.New(repo, cache)
	jwtUtil := jwt.New[actions.TokenPayload]()
	usecases := actions.New(
		app,
		cache,
		jwtUtil,
	)
	authMiddleware := auth.New(usecases)
	minifyer := minify.New()
	minifyer.AddFuncRegexp(regexp.MustCompile("[/+]json$"), json.Minify)

	emailLoginApi := apis.NewUsernameLoginApi(usecases)
	meApi := apis.NewMeApi(usecases)
	accountApi := apis.NewAccountApi(usecases)
	bloodTestApi := apis.NewBloodTestApi(usecases)
	medicineApi := apis.NewMedicineApi(usecases)
	virusApi := apis.NewVirusApi(usecases)
	addressApi := apis.NewAddressApi(usecases)
	patientApi := apis.NewPatientApi(usecases)
	diagnosisApi := apis.NewDiagnosisApi(usecases)

	v1ApisHandler := http.NewServeMux()
	v1ApisHandler.HandleFunc("POST /login/username", emailLoginApi.HandleUsernameLogin)

	v1ApisHandler.HandleFunc("GET /me/auth", authMiddleware.AuthApi(meApi.HandleAuthCheck))
	v1ApisHandler.HandleFunc("GET /me/logout", authMiddleware.AuthApi(meApi.HandleLogout))

	v1ApisHandler.HandleFunc("GET /accounts/{id}", authMiddleware.AuthApi(accountApi.HandleGetAccount))
	v1ApisHandler.HandleFunc("DELETE /accounts/{id}", authMiddleware.AuthApi(accountApi.HandleDeleteAccount))
	v1ApisHandler.HandleFunc("PUT /accounts/{id}", authMiddleware.AuthApi(accountApi.HandleUpdateAccount))
	v1ApisHandler.HandleFunc("POST /accounts/admin", authMiddleware.AuthApi(accountApi.HandleCreateAdminAccount))
	v1ApisHandler.HandleFunc("POST /accounts/secritary", authMiddleware.AuthApi(accountApi.HandleCreateSecritaryAccount))
	v1ApisHandler.HandleFunc("GET /accounts", authMiddleware.AuthApi(accountApi.HandleListAllAccounts))

	v1ApisHandler.HandleFunc("POST /bloodtests", authMiddleware.AuthApi(bloodTestApi.HandleCreateBloodTest))
	v1ApisHandler.HandleFunc("GET /bloodtests/{id}", authMiddleware.AuthApi(bloodTestApi.HandleGetBloodTest))
	v1ApisHandler.HandleFunc("GET /bloodtests/all", authMiddleware.AuthApi(bloodTestApi.HandleListBloodTests))
	v1ApisHandler.HandleFunc("DELETE /bloodtests/{id}", authMiddleware.AuthApi(bloodTestApi.HandleDeleteBloodTest))

	v1ApisHandler.HandleFunc("POST /diagnoses", authMiddleware.AuthApi(diagnosisApi.HandleCreateDiagnosis))
	v1ApisHandler.HandleFunc("GET /diagnoses", authMiddleware.AuthApi(diagnosisApi.HandleListDiagnosiss))
	v1ApisHandler.HandleFunc("DELETE /diagnoses/{id}", authMiddleware.AuthApi(diagnosisApi.HandleDeleteDiagnosis))

	v1ApisHandler.HandleFunc("POST /viruses", authMiddleware.AuthApi(virusApi.HandleCreateVirus))
	v1ApisHandler.HandleFunc("GET /viruses", authMiddleware.AuthApi(virusApi.HandleListViruses))
	v1ApisHandler.HandleFunc("DELETE /viruses/{id}", authMiddleware.AuthApi(virusApi.HandleDeleteVirus))

	v1ApisHandler.HandleFunc("POST /medicines", authMiddleware.AuthApi(medicineApi.HandleCreateMedicine))
	v1ApisHandler.HandleFunc("GET /medicines/all", authMiddleware.AuthApi(medicineApi.HandleListMedicines))
	v1ApisHandler.HandleFunc("GET /medicines/{id}", authMiddleware.AuthApi(medicineApi.HandleGetMedicine))
	v1ApisHandler.HandleFunc("PUT /medicines/{id}/amount", authMiddleware.AuthApi(medicineApi.HandleUpdateMedicineAmount))
	v1ApisHandler.HandleFunc("DELETE /medicines/{id}", authMiddleware.AuthApi(medicineApi.HandleDeleteMedicine))

	v1ApisHandler.HandleFunc(
		"GET /addresses/goveronate/{goveronate}/suburb/{suburb}/street/{street}",
		authMiddleware.AuthApi(addressApi.HandleFindAddress))

	v1ApisHandler.HandleFunc("POST /patients", authMiddleware.AuthApi(patientApi.HandleCreatePatient))
	v1ApisHandler.HandleFunc("GET /patients/{id}/card", authMiddleware.AuthApi(patientApi.HandleGenerateCard))
	v1ApisHandler.HandleFunc("DELETE /patients/{id}", authMiddleware.AuthApi(patientApi.HandleDeletePatient))
	v1ApisHandler.HandleFunc("GET /patients/{id}", authMiddleware.AuthApi(patientApi.HandleGetPatient))
	v1ApisHandler.HandleFunc("GET /patients/last", authMiddleware.AuthApi(patientApi.HandleListLastPatients))
	v1ApisHandler.HandleFunc(
		"GET /patients/public-id/{public_id}/first-name/{first_name}/last-name/{last_name}/father-name/{father_name}/mother-name/{mother_name}/national-id/{national_id}/phone-number/{phone_number}",
		authMiddleware.AuthApi(patientApi.HandleFindPatients))

	v1ApisHandler.HandleFunc("POST /patients/bloodtest", authMiddleware.AuthApi(patientApi.HandleCreatePatientBloodTestResult))
	v1ApisHandler.HandleFunc("PUT /patients/{id}/bloodtest/{btr_id}/pending", authMiddleware.AuthApi(patientApi.HandleUpdatePendingBloodTestResult))
	v1ApisHandler.HandleFunc("POST /patients/{id}/checkup", authMiddleware.AuthApi(patientApi.HandleCheckUp))
	v1ApisHandler.HandleFunc("POST /patients/diagnosis", authMiddleware.AuthApi(patientApi.HandleCreatePatientDiagnosisResult))
	v1ApisHandler.HandleFunc("POST /patients/{id}/joints-evaluation", authMiddleware.AuthApi(patientApi.HandleCreatePatientJointsEvaluation))
	v1ApisHandler.HandleFunc("GET /patients/{id}/joints-evaluations", authMiddleware.AuthApi(patientApi.HandleListPatientJointsEvaluations))
	v1ApisHandler.HandleFunc("GET /patients/{id}/visits", authMiddleware.AuthApi(patientApi.HandleListPatientVisits))
	v1ApisHandler.HandleFunc("POST /patients/visit/{visit_id}/medicine/{med_id}", authMiddleware.AuthApi(patientApi.HandleUsePrescribedMedicineForVisit))

	v1ApisHandler.HandleFunc("GET /me/patient/last-visit", authMiddleware.AuthApi(patientApi.HandleGetPatientLastVisit))

	if config.Env().GoEnv == config.GoEnvTest || config.Env().GoEnv == config.GoEnvDev {
		v1ApisHandler.HandleFunc("GET /status", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"message": "yeeehaww"}`))
		})

		v1ApisHandler.HandleFunc("POST /tests/reset/db", func(w http.ResponseWriter, r *http.Request) {
			err := repo.DeleteAll()
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(`{"message": "resetting DB failed"}`))
				return
			}

			_ = repo.CreateSuperAdmin()

			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"message": "yeeehaww"}`))
		})

		v1ApisHandler.HandleFunc("POST /tests/reset/cache", func(w http.ResponseWriter, r *http.Request) {
			err := cache.FlushAll()
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(`{"message": "flushing cache failed"}`))
				return
			}

			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"message": "yeeehaww"}`))
		})
	}

	applicationHandler := http.NewServeMux()
	applicationHandler.Handle("/v1/", http.StripPrefix("/v1", contenttype.Json(v1ApisHandler)))

	log.Info("Starting http server at port " + config.Env().Port)
	switch config.Env().GoEnv {
	case config.GoEnvBeta, config.GoEnvDev, config.GoEnvTest:
		log.Fatalln(http.ListenAndServe(":"+config.Env().Port, logger.Handler(applicationHandler)))
	case config.GoEnvProd:
		log.Fatalln(http.ListenAndServe(":"+config.Env().Port, minifyer.Middleware(applicationHandler)))
	}
}
