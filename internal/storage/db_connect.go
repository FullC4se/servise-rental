package storage

import (
	"context"
	"database/sql"
	"fmt"
)

func DbConnect(ctx context.Context, driverName, dataSourceName string) (*sql.DB, error) {
	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("Ошибка подключения к базе данных!\n\nПроверьте следующее:\n\t- Драйвер базы данных\n\t- Правильность строки подключения\n\nОшибка: %v", err)
	}
	err = db.PingContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("Ошибка подключения к базе данных!\n\nПроверьте следующее:\n\t- Логин\n\t- Пароль\n\t- Хост\n\t- Порт\n\nОшибка: %v", err)
	}
	return db, nil
}
