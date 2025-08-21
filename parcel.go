package main

import (
	"database/sql"
	"log/slog"
)

type ParcelStore struct {
	db *sql.DB
}

func NewParcelStore(db *sql.DB) ParcelStore {
	return ParcelStore{db: db}
}

func (s ParcelStore) Add(p Parcel) (int, error) {
	// реализуйте добавление строки в таблицу parcel, используйте данные из переменной p
	res, err := s.db.Exec(
		`INSERT INTO parcel (client, status, address, created_at)
	VALUES (:Client, :Status, :Address, :Created_at)`,
		sql.Named("Client", p.Client),
		sql.Named("Status", p.Status),
		sql.Named("Address", p.Address),
		sql.Named("Created_at", p.CreatedAt),
	)
	if err != nil {
		return 0, err
	}
	// верните идентификатор последней добавленной записи
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	slog.Info("добавлена новая посылка", "id", id, "number", p.Number)
	return int(id), nil
}

func (s ParcelStore) Get(number int) (Parcel, error) {
	// реализуйте чтение строки по заданному number
	// здесь из таблицы должна вернуться только одна строка
	row := s.db.QueryRow(`SELECT number, client, status, address, created_at 
	                    FROM parcel WHERE number = :number`,
		sql.Named("number", number))
	// заполните объект Parcel данными из таблицы
	p := Parcel{}

	err := row.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
	if err != nil {
		slog.Error("не удалось выполнить запрос Get", "err", err, "по номеру", number)
		return Parcel{}, err
	}

	slog.Info("получена посылка", "number", number)

	return p, nil
}

func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {
	// реализуйте чтение строк из таблицы parcel по заданному client
	// здесь из таблицы может вернуться несколько строк
	rows, err := s.db.Query(`SELECT number, client, status, address, created_at
	                      FROM parcel WHERE client = :client`,
		sql.Named("client", client))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// заполните срез Parcel данными из таблицы
	var res []Parcel
	for rows.Next() {
		var p Parcel
		if err := rows.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt); err != nil {
			return nil, err
		}
		res = append(res, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	slog.Info("получены посылки по клиенту", "client", client)

	return res, nil
}

func (s ParcelStore) SetStatus(number int, status string) error {
	// реализуйте обновление статуса в таблице parcel
	res, err := s.db.Exec(`UPDATE parcel
	                     SET status = :status 
											 WHERE number = :number`,
		sql.Named("status", status),
		sql.Named("number", number))
	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	slog.Info("обновлён статус посылки", "number", number, "rows", affected)

	return nil
}

func (s ParcelStore) SetAddress(number int, address string) error {
	// реализуйте обновление адреса в таблице parcel
	// менять адрес можно только если значение статуса registered
	
	res, err := s.db.Exec(`UPDATE parcel
	                       SET address = :address
												 WHERE number = :number AND status = :status`,
		sql.Named("address", address),
		sql.Named("number", number),
		sql.Named("status", ParcelStatusRegistered))
	if err != nil {
		return err
	}

	affected, _ := res.RowsAffected()
	slog.Info("обновлён адрес посылки", "number", number, "rows", affected)

	return nil
}

func (s ParcelStore) Delete(number int) error {
	// реализуйте удаление строки из таблицы parcel
	// удалять строку можно только если значение статуса registered

	res, err := s.db.Exec(
		`DELETE FROM parcel WHERE number = :number AND status = :status`,
		sql.Named("number", number),
		sql.Named("status", ParcelStatusRegistered),
	)
	if err != nil {
		return err
	}

	affected, _ := res.RowsAffected()
	slog.Info("посылка удалена", "number", number, "rows", affected)

	return nil
}
