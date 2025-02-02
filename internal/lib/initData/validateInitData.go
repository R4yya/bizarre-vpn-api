package initData

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

type InitDataUserPayload struct {
	TelegramID      int64  `json:"id"`
	FirstName       string `json:"first_name"`
	LastName        string `json:"last_name"`
	Username        string `json:"username"`
	LanguageCode    string `json:"language_code"`
	AllowsWriteToPm bool   `json:"allows_write_to_pm"`
	PhotoUrl        string `json:"photo_url"`
}

func ValidateInitData(
	initDataStr string,
	botToken string,
	expirationTime time.Duration,
) (payload *InitDataUserPayload, isOk bool, Err error) {
	op := "internal.lib.initData.validateInitData"

	// Парсим initData в map
	initData, err := url.ParseQuery(initDataStr)
	if err != nil {
		mes := "error when parsing initData"
		return nil, false, fmt.Errorf("%v: %v: %w", op, mes, err)
	}

	// извлекаем hash
	receivedHash := initData.Get("hash")
	if receivedHash == "" {
		mes := "hash is not defined"
		return nil, false, fmt.Errorf("%v: %v", op, mes)
	}

	// Удаляем хэш из данных для проверки
	initData.Del("hash")

	// Сортируем ключи и собираем строку для хэширования
	var keys []string
	for key := range initData {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	var dataCheckStringBuilder strings.Builder
	for _, key := range keys {
		dataCheckStringBuilder.WriteString(key)
		dataCheckStringBuilder.WriteString("=")
		dataCheckStringBuilder.WriteString(initData.Get(key))
		dataCheckStringBuilder.WriteString("\n")
	}
	dataCheckString := strings.TrimSuffix(dataCheckStringBuilder.String(), "\n")

	// Вычисляем HMAC-SHA-256
	secretKey := hmac.New(sha256.New, []byte("WebAppData"))
	secretKey.Write([]byte(botToken))
	computedHash := hmac.New(sha256.New, secretKey.Sum(nil))
	computedHash.Write([]byte(dataCheckString))
	expectedHash := hex.EncodeToString(computedHash.Sum(nil))

	// Сравниваем хэши
	if receivedHash != expectedHash {
		mes := "wrong hash"
		return nil, false, fmt.Errorf("%v: %v", op, mes)
	}

	// Проверяем просроченность initData
	authDateStr := initData.Get("auth_date")
	if authDateStr == "" {
		mes := "not found "
		return nil, false, fmt.Errorf("%v: %v", op, mes)
	}

	authDate, err := strconv.ParseInt(authDateStr, 10, 64)
	if err != nil {
		mes := "error when parsing auth_date"
		return nil, false, fmt.Errorf("%v: %v: %w", op, mes, err)
	}

	authTime := time.Unix(authDate, 0)
	if time.Since(authTime) > expirationTime {
		return nil, false, fmt.Errorf("initData is expire")
	}

	userJsonString := initData.Get("user")

	var initDataUserPayload InitDataUserPayload

	err = json.Unmarshal([]byte(userJsonString), &payload)

	if err != nil {
		mes := "error when unmarshal initDataUserPayload"
		return nil, false, fmt.Errorf("%v: %v: %w", op, mes, err)
	}

	return &initDataUserPayload, true, nil
}
