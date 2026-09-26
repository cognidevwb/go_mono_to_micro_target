# Performance

Connection budget: 100 connections per database server, divided across each service's replicas (`DB_MAX_CONNS`). Cross-service calls carry a 5 s timeout. Scale a hot service alone — its HPA is independent.
