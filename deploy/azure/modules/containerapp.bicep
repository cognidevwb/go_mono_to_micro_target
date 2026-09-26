param name string
param location string
param environmentId string
param registry string
param identityId string
param external bool
param image string = '${registry}/${name}:latest'

resource app 'Microsoft.App/containerApps@2024-03-01' = {
  name: name
  location: location
  identity: { type: 'UserAssigned', userAssignedIdentities: { '${identityId}': {} } }
  properties: {
    managedEnvironmentId: environmentId
    configuration: {
      ingress: { external: external, targetPort: 8080 }
      registries: [ { server: registry, identity: identityId } ]
    }
    template: {
      containers: [
        {
          name: name
          image: image
          resources: { cpu: json('0.5'), memory: '1Gi' }
          probes: [
            { type: 'Liveness', httpGet: { path: '/healthz', port: 8080 } }
            { type: 'Readiness', httpGet: { path: '/readyz', port: 8080 } }
          ]
        }
      ]
      scale: { minReplicas: 1, maxReplicas: 5 }
    }
  }
}

output fqdn string = app.properties.configuration.ingress.fqdn
