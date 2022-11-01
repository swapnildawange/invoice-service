package dl

import "context"

type Storer interface {
	CreateInvoice(ctx context.Context)
}
