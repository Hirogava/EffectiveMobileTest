package api

// @title Effective Mobile API
// @version 1.0
// @description API для управления данными о людях
// @host localhost:8080
// @BasePath /api

// @Summary Добавить нового человека
// @Description Добавляет нового человека в базу данных с автоматическим определением возраста, пола и национальности
// @Tags people
// @Accept json
// @Produce json
// @Param person body models.PeopleFromRequest true "Данные человека"
// @Success 200 {object} models.People
// @Failure 400 {object} map[string]string
// @Failure 415 {object} map[string]string
// @Failure 503 {object} map[string]string
// @Router /people [post]
func SwaggerAddPeople() {}

// @Summary Получить данные о человеке
// @Description Возвращает полные данные о человеке по ID
// @Tags people
// @Produce json
// @Param id path int true "ID человека"
// @Success 200 {object} models.People
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /people/{id} [get]
func SwaggerGetPeople() {}

// @Summary Удалить человека
// @Description Удаляет человека из базы данных по ID
// @Tags people
// @Produce json
// @Param id path int true "ID человека"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /people/{id} [delete]
func SwaggerDeletePeople() {}

// @Summary Обновить данные человека
// @Description Обновляет данные о человеке. При изменении имени автоматически обновляет возраст, пол и национальность
// @Tags people
// @Accept json
// @Produce json
// @Param id path int true "ID человека"
// @Param person body models.PeopleFromRequest true "Новые данные человека"
// @Success 200 {object} models.People
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /people/{id} [put]
func SwaggerUpdatePeople() {}

// @Summary Получить список людей
// @Description Возвращает список людей с возможностью фильтрации по возрасту, полу и стране
// @Tags people
// @Produce json
// @Param age_from query int false "Минимальный возраст"
// @Param age_to query int false "Максимальный возраст"
// @Param gender query string false "Пол (male/female)"
// @Param country query string false "Код страны"
// @Param page query int false "Номер страницы" default(1)
// @Param limit query int false "Количество записей на странице" default(10)
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /people [get]
func SwaggerGetAllPeople() {}
