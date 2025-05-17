package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Структура для парсинга ответа Weatherstack
type WeatherstackResponse struct {
	Location struct {
		Name string `json:"name"`
	} `json:"location"`
	Current struct {
		Temperature         int      `json:"temperature"`          // Температура в °C
		WeatherDescriptions []string `json:"weather_descriptions"` // Описание погоды
	} `json:"current"`
}

func GetWeather(city, apiKey string) (string, error) {
	url := fmt.Sprintf("http://api.weatherstack.com/current?access_key=%s&query=%s&units=m", apiKey, city)
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("ошибка запроса: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API вернуло статус: %d", resp.StatusCode)
	}

	var data WeatherstackResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", fmt.Errorf("ошибка парсинга JSON: %v", err)
	}

	if len(data.Current.WeatherDescriptions) == 0 {
		return "", fmt.Errorf("город не найден")
	}

	description := data.Current.WeatherDescriptions[0]
	return fmt.Sprintf("Погода в %s: %d°C (%s)",
		data.Location.Name,
		data.Current.Temperature,
		description,
	), nil
}

var ctx = context.Background()
var rdb *redis.Client

// Инициализация Redis
func initRedis() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Адрес Redis
		Password: "",               // Пароль (если есть)
		DB:       0,                // Номер базы данных
	})
}

// Получение погоды с кешированием
func GetWeather(city, apiKey string) (string, error) {
	// Проверяем кеш
	cachedWeather, err := rdb.Get(ctx, "weather:"+city).Result()
	if err == nil {
		return cachedWeather, nil // Возвращаем данные из кеша
	}

	// Если в кеше нет, делаем запрос к API
	url := fmt.Sprintf("http://api.weatherstack.com/current?access_key=%s&query=%s&units=m", apiKey, city)
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("ошибка запроса: %v", err)
	}
	defer resp.Body.Close()

	var data WeatherstackResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", fmt.Errorf("ошибка парсинга JSON: %v", err)
	}

	if len(data.Current.WeatherDescriptions) == 0 {
		return "", fmt.Errorf("город не найден")
	}

	description := data.Current.WeatherDescriptions[0]
	weatherText := fmt.Sprintf("Погода в %s: %d°C (%s)", data.Location.Name, data.Current.Temperature, description)

	// Сохраняем в кеш на 30 минут
	err = rdb.Set(ctx, "weather:"+city, weatherText, 30*time.Minute).Err()
	if err != nil {
		return weatherText, fmt.Errorf("не удалось сохранить в кеш: %v", err)
	}

	return weatherText, nil
}

var ctx = context.Background()
var rdb *redis.Client

// Инициализация Redis
func initRedis() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Адрес Redis
		Password: "",               // Пароль (если есть)
		DB:       0,                // Номер базы данных
	})
}

// Получение погоды с кешированием
func GetWeather(city, apiKey string) (string, error) {
	// Проверяем кеш
	cachedWeather, err := rdb.Get(ctx, "weather:"+city).Result()
	if err == nil {
		return cachedWeather, nil // Возвращаем данные из кеша
	}

	// Если в кеше нет, делаем запрос к API
	url := fmt.Sprintf("http://api.weatherstack.com/current?access_key=%s&query=%s&units=m", apiKey, city)
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("ошибка запроса: %v", err)
	}
	defer resp.Body.Close()

	var data WeatherstackResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", fmt.Errorf("ошибка парсинга JSON: %v", err)
	}

	if len(data.Current.WeatherDescriptions) == 0 {
		return "", fmt.Errorf("город не найден")
	}

	description := data.Current.WeatherDescriptions[0]
	weatherText := fmt.Sprintf("Погода в %s: %d°C (%s)", data.Location.Name, data.Current.Temperature, description)

	// Сохраняем в кеш на 30 минут
	err = rdb.Set(ctx, "weather:"+city, weatherText, 30*time.Minute).Err()
	if err != nil {
		return weatherText, fmt.Errorf("не удалось сохранить в кеш: %v", err)
	}

	return weatherText, nil
}
