// Azure Container Apps landing zone for the fleet: Log Analytics, a Container
// Apps environment, ACR, one user-assigned identity (AcrPull) and one container
// app per service. Deploy with `azd up` (azure.yaml) or
// `az deployment group create -f deploy/azure/main.bicep -p deploy/azure/main.parameters.json`.
targetScope = 'resourceGroup'

param location string = resourceGroup().location
param prefix string = 'platform'

resource logs 'Microsoft.OperationalInsights/workspaces@2023-09-01' = {
  name: '${prefix}-logs'
  location: location
  properties: { sku: { name: 'PerGB2018' }, retentionInDays: 30 }
}

resource env 'Microsoft.App/managedEnvironments@2024-03-01' = {
  name: '${prefix}-env'
  location: location
  properties: {
    appLogsConfiguration: {
      destination: 'log-analytics'
      logAnalyticsConfiguration: { customerId: logs.properties.customerId, sharedKey: logs.listKeys().primarySharedKey }
    }
  }
}

resource acr 'Microsoft.ContainerRegistry/registries@2023-07-01' = {
  name: replace('${prefix}acr', '-', '')
  location: location
  sku: { name: 'Basic' }
  properties: { adminUserEnabled: false }
}

resource identity 'Microsoft.ManagedIdentity/userAssignedIdentities@2023-01-31' = {
  name: '${prefix}-apps'
  location: location
}

resource acrPull 'Microsoft.Authorization/roleAssignments@2022-04-01' = {
  name: guid(acr.id, identity.id, 'acrpull')
  scope: acr
  properties: {
    principalId: identity.properties.principalId
    principalType: 'ServicePrincipal'
    roleDefinitionId: subscriptionResourceId('Microsoft.Authorization/roleDefinitions', '7f951dda-4ed3-4680-a7ca-43fe172d538d')
  }
}

module catalog_service 'modules/containerapp.bicep' = {
  name: 'catalog-service'
  params: {
    name: 'catalog-service'
    location: location
    environmentId: env.id
    registry: acr.properties.loginServer
    identityId: identity.id
    external: false
  }
}

module customers_service 'modules/containerapp.bicep' = {
  name: 'customers-service'
  params: {
    name: 'customers-service'
    location: location
    environmentId: env.id
    registry: acr.properties.loginServer
    identityId: identity.id
    external: false
  }
}

module inventory_service 'modules/containerapp.bicep' = {
  name: 'inventory-service'
  params: {
    name: 'inventory-service'
    location: location
    environmentId: env.id
    registry: acr.properties.loginServer
    identityId: identity.id
    external: false
  }
}

module payments_service 'modules/containerapp.bicep' = {
  name: 'payments-service'
  params: {
    name: 'payments-service'
    location: location
    environmentId: env.id
    registry: acr.properties.loginServer
    identityId: identity.id
    external: false
  }
}

module orders_service 'modules/containerapp.bicep' = {
  name: 'orders-service'
  params: {
    name: 'orders-service'
    location: location
    environmentId: env.id
    registry: acr.properties.loginServer
    identityId: identity.id
    external: false
  }
}

module gateway 'modules/containerapp.bicep' = {
  name: 'gateway'
  params: {
    name: 'gateway'
    location: location
    environmentId: env.id
    registry: acr.properties.loginServer
    identityId: identity.id
    external: true
  }
}
