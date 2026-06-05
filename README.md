# loyalty-core-backoffice

Core del backoffice de la loyalty platform.

## Propósito
Modelar capacidades operativas internas propias del backoffice.

## Integraciones esperadas
- consumo complementario de `loyalty-core-points`

## Estado
Bootstrap funcional con contrato mock mínimo para capacidades operativas y alertas.

## Endpoints
- `GET /health`
- `GET /metrics`
- `GET /v1/backoffice-capabilities`
- `GET /v1/operational-alerts`

## Notas técnicas
- Implementación Go stdlib + Prometheus
- Sin persistencia todavía; usa memoria/mock para bootstrap
- Pensado para evolucionar luego hacia agregación real con `loyalty-core-points`
