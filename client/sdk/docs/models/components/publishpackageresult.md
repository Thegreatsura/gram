# PublishPackageResult

## Example Usage

```typescript
import { PublishPackageResult } from "@gram/client/models/components";

let value: PublishPackageResult = {
  package: {
    createdAt: new Date("2023-03-12T09:33:18.947Z"),
    id: "<id>",
    name: "<value>",
    organizationId: "<id>",
    projectId: "<id>",
    updatedAt: new Date("2025-06-16T13:02:48.836Z"),
  },
  version: {
    createdAt: new Date("2023-12-27T12:10:52.007Z"),
    deploymentId: "<id>",
    id: "<id>",
    packageId: "<id>",
    semver: "<value>",
    visibility: "<value>",
  },
};
```

## Fields

| Field                                                                  | Type                                                                   | Required                                                               | Description                                                            |
| ---------------------------------------------------------------------- | ---------------------------------------------------------------------- | ---------------------------------------------------------------------- | ---------------------------------------------------------------------- |
| `package`                                                              | [components.Package](../../models/components/package.md)               | :heavy_check_mark:                                                     | N/A                                                                    |
| `version`                                                              | [components.PackageVersion](../../models/components/packageversion.md) | :heavy_check_mark:                                                     | N/A                                                                    |