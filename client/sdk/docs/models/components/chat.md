# Chat

## Example Usage

```typescript
import { Chat } from "@gram/client/models/components";

let value: Chat = {
  createdAt: new Date("2023-08-22T03:55:40.546Z"),
  id: "<id>",
  messages: [
    {
      createdAt: new Date("2024-01-23T04:53:42.435Z"),
      id: "<id>",
      model: "Durango",
      role: "<value>",
    },
  ],
  numMessages: 338963,
  title: "<value>",
  updatedAt: new Date("2024-06-25T13:21:36.855Z"),
  userId: "<id>",
};
```

## Fields

| Field                                                                                         | Type                                                                                          | Required                                                                                      | Description                                                                                   |
| --------------------------------------------------------------------------------------------- | --------------------------------------------------------------------------------------------- | --------------------------------------------------------------------------------------------- | --------------------------------------------------------------------------------------------- |
| `createdAt`                                                                                   | [Date](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Date) | :heavy_check_mark:                                                                            | When the chat was created.                                                                    |
| `id`                                                                                          | *string*                                                                                      | :heavy_check_mark:                                                                            | The ID of the chat                                                                            |
| `messages`                                                                                    | [components.ChatMessage](../../models/components/chatmessage.md)[]                            | :heavy_check_mark:                                                                            | The list of messages in the chat                                                              |
| `numMessages`                                                                                 | *number*                                                                                      | :heavy_check_mark:                                                                            | The number of messages in the chat                                                            |
| `title`                                                                                       | *string*                                                                                      | :heavy_check_mark:                                                                            | The title of the chat                                                                         |
| `updatedAt`                                                                                   | [Date](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Date) | :heavy_check_mark:                                                                            | When the chat was last updated.                                                               |
| `userId`                                                                                      | *string*                                                                                      | :heavy_check_mark:                                                                            | The ID of the user who created the chat                                                       |