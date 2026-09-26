package proxy

// Routes is the prefix table, one row per extracted service's path prefix.
// Regenerated from the plan on every scaffold run.
func Routes() []Route {
	return []Route{
		{Prefix: "/api/products", Service: "catalog-service", Upstream: "http://catalog-service:8080"},
		{Prefix: "/api/customers", Service: "customers-service", Upstream: "http://customers-service:8080"},
		{Prefix: "/api/inventory", Service: "inventory-service", Upstream: "http://inventory-service:8080"},
		{Prefix: "/api/payments", Service: "payments-service", Upstream: "http://payments-service:8080"},
		{Prefix: "/api/orders", Service: "orders-service", Upstream: "http://orders-service:8080"},
	}
}
