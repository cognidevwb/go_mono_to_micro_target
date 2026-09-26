# 5. Strangler cutover order

Status: accepted

Migration approach: strangler. The gateway routes each extracted prefix to its service and everything else to the monolith (`LEGACY_UPSTREAM`); leaves move first, the saga orchestrator last.
