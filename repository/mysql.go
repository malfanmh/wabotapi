package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/malfanmh/wabotapi/model"
	"time"
)

type mysqlRepository struct {
	db *sqlx.DB
}

func NewMysql(db *sqlx.DB) *mysqlRepository {
	return &mysqlRepository{db}
}

func (r *mysqlRepository) GetClientByWAPhoneID(ctx context.Context, waPhoneID string) (result model.Client, err error) {
	q := `SELECT id, name, hash, token, wa_phone, wa_phone_id FROM clients WHERE wa_phone_id = ?`
	err = r.db.GetContext(ctx, &result, r.db.Rebind(q), waPhoneID)
	return
}

func (r *mysqlRepository) GetMessageFlow(ctx context.Context, clientID int64, access model.Access, keyword string, seq string, limit int) (result []model.MessageFlow, err error) {
	var args = []interface{}{keyword, clientID, access}
	q := `select keyword,mf.message_id, mf.access,m.type,mf.validate_input,mf.seq, m.slug, mf.checkout
			from message_flows mf
				 inner join messages m on mf.message_id = m.id
		where mf.display = 1
		  and keyword = ? 
		  and mf.client_id = ?
		  and mf.access = ?		   
		`
	if seq != "" {
		q += " and mf.seq >= ? "
		args = append(args, seq)
	}
	q += `order by mf.seq`
	if limit > 0 {
		q += fmt.Sprintf(" limit %d", limit)
	}
	fmt.Println(q)
	fmt.Println(args)
	err = r.db.SelectContext(ctx, &result, r.db.Rebind(q), args...)
	return
}

func (r *mysqlRepository) GetMessageFlowBySlug(ctx context.Context, clientID int64, slug string) (flow model.MessageFlow, err error) {
	q := `select keyword,mf.message_id, mf.access,m.type,mf.validate_input,mf.seq, m.slug,mf.checkout
			from messages m inner join message_flows mf on m.id = mf.message_id
			where m.slug = ?
			  and m.client_id = ?`
	err = r.db.GetContext(ctx, &flow, r.db.Rebind(q), slug, clientID)
	return
}

func (r *mysqlRepository) GetNextFlow(ctx context.Context, clientID int64, access model.Access, keyword string, seq string) (result model.MessageFlow, err error) {
	var args = []interface{}{keyword, clientID, access, seq}
	q := `select keyword,mf.message_id, mf.access,m.type,mf.validate_input,mf.seq, m.slug,mf.checkout
			from message_flows mf
				 inner join messages m on mf.message_id = m.id
		where mf.display = 1
		  and keyword = ? 
		  and mf.client_id = ?
		  and mf.access = ?	
		  and seq > ?
		  order by mf.seq
		  limit 1
		  `
	err = r.db.GetContext(ctx, &result, r.db.Rebind(q), args...)
	return
}

func (r *mysqlRepository) GetMessage(ctx context.Context, clientID, messageID int64) (msg model.Message, err error) {
	q := `select slug, type, button, header_text, preview_url, body_text, footer_text, with_metadata
			from messages
			where client_id = ?
			and id = ?`
	err = r.db.GetContext(ctx, &msg, r.db.Rebind(q), clientID, messageID)
	return
}

func (r *mysqlRepository) GetMessageBySlug(ctx context.Context, clientID int64, slug string) (msg model.Message, err error) {
	q := `select id, slug, type, button, header_text, preview_url, body_text, footer_text, with_metadata
			from messages
			where client_id = ?
			and slug = ?`
	err = r.db.GetContext(ctx, &msg, r.db.Rebind(q), clientID, slug)
	return
}

func (r *mysqlRepository) GetMessageAction(ctx context.Context, messageID int64, access model.Access) (result []model.MessageAction, err error) {
	q := `select slug,title,description
			from message_actions
			where message_id = ?
			 and access = ?
			 and display = 1
			order by seq;`
	err = r.db.SelectContext(ctx, &result, r.db.Rebind(q), messageID, access)
	return
}

func (r *mysqlRepository) GetCustomerByWAID(ctx context.Context, clientID int64, waid string) (result model.Customer, err error) {
	q := `SELECT id, client_id, wa_id, email, full_name, address, birth_date, identity_number, identity_type, gender, access, created_at, updated_at FROM customers WHERE client_id = ? and wa_id = ?`
	err = r.db.GetContext(ctx, &result, r.db.Rebind(q), clientID, waid)
	return
}

func (r *mysqlRepository) GetSession(ctx context.Context, clientID int64, waid, input string) (result model.Session, err error) {
	q := `SELECT id, client_id, wa_id, access, seq, slug, input, created_at, expired_at FROM session WHERE client_id = ? and wa_id = ?`
	err = r.db.GetContext(ctx, &result, r.db.Rebind(q), clientID, waid)
	if errors.Is(err, sql.ErrNoRows) {
		now := time.Now()
		expired := now.Add(24 * time.Hour)
		q := `INSERT INTO session(client_id, wa_id, access, seq, input, expired_at) VALUES (?,?,?,?,?,?)`
		res, errI := r.db.ExecContext(ctx, r.db.Rebind(q), clientID, waid, 0, "", input, expired.Format(time.RFC3339))
		if errI != nil {
			err = errI
			return
		}
		id, _ := res.LastInsertId()
		result = model.Session{
			ID:        id,
			ClientID:  clientID,
			WAID:      waid,
			Access:    0,
			Seq:       "",
			Input:     input,
			CreatedAt: now.Format(time.RFC3339),
			ExpiredAt: expired.Format(time.RFC3339),
		}

	}
	if err != nil {
		return result, err
	}
	return
}

func (r *mysqlRepository) UpdateSession(ctx context.Context, session model.Session) (err error) {
	expired := time.Now().Add(24 * time.Hour)
	q := `UPDATE session SET seq = ? ,slug = ?, input = ?,access = ?, expired_at = ? WHERE client_id = ? and wa_id = ?`
	_, err = r.db.ExecContext(ctx, r.db.Rebind(q), session.Seq, session.Slug, session.Input, session.Access, expired.Format(time.RFC3339), session.ClientID, session.WAID)
	return
}

func (r *mysqlRepository) InsertCustomer(ctx context.Context, customer model.Customer) (err error) {
	q := `INSERT INTO customers (client_id,wa_id, access) VALUES (?,?,?)`
	_, err = r.db.ExecContext(ctx, r.db.Rebind(q), customer.ClientID, customer.WAID, model.AccessNew)
	return
}

func (r *mysqlRepository) UpdateCustomer(ctx context.Context, customer model.Customer) (err error) {
	builder := sq.Update("customers")
	if customer.Email.Valid {
		builder = builder.Set("email", customer.Email.V)
	}

	if customer.FullName.Valid {
		builder = builder.Set("full_name", customer.FullName.V)
	}

	if customer.BirthDate.Valid {
		birtDate := customer.BirthDate.V
		if len(birtDate) > 10 {
			if t, errP := time.Parse(time.RFC3339, birtDate); errP == nil {
				birtDate = t.Format(time.DateOnly)
			}
		}
		builder = builder.Set("birth_date", birtDate)
	}
	if customer.IdentityNumber.Valid {
		builder = builder.Set("identity_number", customer.IdentityNumber.V)
		builder = builder.Set("identity_type", "ktp")
	}
	if customer.Address.Valid {
		builder = builder.Set("address", customer.Address.V)
	}
	if customer.Gender.Valid {
		builder = builder.Set("gender", customer.Gender.V)
	}
	if customer.Access.Valid {
		builder = builder.Set("status", customer.Access.V.Int())
	}
	builder = builder.Where(sq.Eq{
		"client_id": customer.ClientID,
		"wa_id":     customer.WAID,
	})

	q, args, err := builder.ToSql()
	fmt.Println("update customer", q, args, err)

	if err != nil {
		return
	}

	_, err = r.db.ExecContext(ctx, q, args...)
	return
}

func (r *mysqlRepository) CreatePayment(ctx context.Context, data model.Payment) (lastID int64, err error) {
	q := `insert into payments(id, customer_id, ref_id, payment_provider, payment_ref_id, payment_type, payment_code,payment_item, status, amount, fee)
				values (:id,  :customer_id, :ref_id, :payment_provider, :payment_ref_id,:payment_type, :payment_code, :payment_item, :status, :amount, :fee)`

	result, err := r.db.NamedExecContext(ctx, q, data)
	if err != nil {
		return 0, err
	}
	lastID, err = result.LastInsertId()
	return
}

func (r *mysqlRepository) GetPaymentCustomer(ctx context.Context, id string) (result model.PaymentCustomer, err error) {
	q := `select c.wa_id,cl.wa_phone_id, customer_id,t.id,c.client_id,c.full_name
			from payments t 
			    inner join customers c on t.customer_id = c.id
			    inner join clients cl on c.client_id = cl.id
			where t.id = ?`

	r.db.GetContext(ctx, &result, r.db.Rebind(q), id)
	return
}

func (r *mysqlRepository) UpdatePayment(ctx context.Context, data model.Payment) (err error) {
	builder := sq.Update("payments")
	if data.Status.Valid {
		builder = builder.Set("status", data.Status.V)
	}
	if data.RefID.Valid {
		builder = builder.Set("ref_id", data.RefID.V)
	}
	if data.PaymentProvider.Valid {
		builder = builder.Set("payment_provider", data.PaymentProvider.V)
	}
	if data.PaymentRefID.Valid {
		builder = builder.Set("payment_ref_id", data.PaymentRefID.V)
	}
	if data.PaymentType.Valid {
		builder = builder.Set("payment_type", data.PaymentType.V)
	}
	if data.PaymentCode.Valid {
		builder = builder.Set("payment_code", data.PaymentCode.V)
	}
	if data.PaymentItem.Valid {
		builder = builder.Set("payment_item", data.PaymentItem.V)
	}
	if data.Amount.Valid {
		builder = builder.Set("amount", data.Amount.Decimal)
	}
	if data.Fee.Valid {
		builder = builder.Set("fee", data.Fee.Decimal)
	}

	if data.ExpiredAt.Valid {
		builder = builder.Set("expired_at", data.ExpiredAt)
	}

	builder = builder.Where(sq.Eq{
		"id": data.ID,
	})

	q, args, err := builder.ToSql()
	fmt.Println("update payments", q, args, err)

	if err != nil {
		return
	}

	_, err = r.db.ExecContext(ctx, q, args...)
	return
}

func (r *mysqlRepository) GetProductBySlug(ctx context.Context, clientID int64, slug string) (result model.Product, err error) {
	q := `select id, name,slug, description, price, stock from products where slug = ?`
	err = r.db.GetContext(ctx, &result, r.db.Rebind(q), slug)
	return
}

func (r *mysqlRepository) GetExpiredPayment(ctx context.Context) (result []model.Payment, err error) {
	q := `select id,status from payments where expired_at < now() and status = 'PENDING';`
	err = r.db.SelectContext(ctx, &result, q)
	return
}
