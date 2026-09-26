# 1. Database per service

Status: accepted

Each service owns its tables in its own PostgreSQL database. The monolith's shared `*gorm.DB` is split; a service reads another's data only through its API or its events.
