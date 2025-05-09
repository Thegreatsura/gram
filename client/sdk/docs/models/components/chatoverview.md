# ChatOverview

## Example Usage

```typescript
import { ChatOverview } from "@gram/client/models/components";

let value: ChatOverview = {
  createdAt: new Date("2024-05-18T14:46:31.663Z"),
  id: "<id>",
  numMessages: 349178,
  title: "<value>",
  updatedAt: new Date("2023-01-03T09:29:42.548Z"),
  userId: "<id>",
};
```

## Fields

| Field                                                                                         | Type                                                                                          | Required                                                                                      | Description                                                                                   |
| --------------------------------------------------------------------------------------------- | --------------------------------------------------------------------------------------------- | --------------------------------------------------------------------------------------------- | --------------------------------------------------------------------------------------------- |
| `createdAt`                                                                                   | [Date](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Date) | :heavy_check_mark:                                                                            | When the chat was created.                                                                    |
| `id`                                                                                          | *string*                                                                                      | :heavy_check_mark:                                                                            | The ID of the chat                                                                            |
| `numMessages`                                                                                 | *number*                                                                                      | :heavy_check_mark:                                                                            | The number of messages in the chat                                                            |
| `title`                                                                                       | *string*                                                                                      | :heavy_check_mark:                                                                            | The title of the chat                                                                         |
| `updatedAt`                                                                                   | [Date](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Date) | :heavy_check_mark:                                                                            | When the chat was last updated.                                                               |
| `userId`                                                                                      | *string*                                                                                      | :heavy_check_mark:                                                                            | The ID of the user who created the chat                                                       |