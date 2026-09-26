# 2. Sagas, not distributed transactions

Status: accepted

A write that spanned contexts in one `db.Transaction` becomes an orchestration saga (`pkg/saga`): ordered steps, each with a compensation, the hardest to undo last. Two-phase commit is not used.
