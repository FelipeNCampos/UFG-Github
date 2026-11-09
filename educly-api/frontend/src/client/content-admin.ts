import type { CancelablePromise } from "./core/CancelablePromise"
import { OpenAPI } from "./core/OpenAPI"
import { request as __request } from "./core/request"

export type TopicBehavior = "atividade" | "explicacao"
export type ActivityType =
  | "MatchWords"
  | "FillBlanks"
  | "Quizz"
  | "PromptTreiner"

export interface ContentPublic {
  id: number
  plan: string
  language: string
  day_id: number
  title: string
  content: string
  created_at?: string | null
}

export interface TopicPublic {
  id: string
  content_id: number
  title?: string | null
  body?: string | null
  behavior: TopicBehavior
  activity_type?: ActivityType | null
  payload: Record<string, unknown>
  order: number
  created_at?: string | null
}

export interface ContentsPublic {
  data: ContentPublic[]
  count: number
}

export interface TopicsPublic {
  data: TopicPublic[]
  count: number
}

export interface MessagePublic {
  message: string
}

export type ContentCreateInput = Omit<ContentPublic, "id" | "created_at">
export type ContentUpdateInput = Partial<ContentCreateInput>
export type TopicCreateInput = Omit<TopicPublic, "id" | "created_at">
export type TopicUpdateInput = Partial<TopicCreateInput>

export class ContentAdminService {
  public static readContents(data: {
    skip?: number
    limit?: number
    plan?: string
    language?: string
    search?: string
  } = {}): CancelablePromise<ContentsPublic> {
    return __request(OpenAPI, {
      method: "GET",
      url: "/api/v1/content/items",
      query: {
        skip: data.skip,
        limit: data.limit,
        plan: data.plan,
        language: data.language,
        search: data.search,
      },
      errors: {
        422: "Validation Error",
      },
    })
  }

  public static createContent(data: {
    requestBody: ContentCreateInput
  }): CancelablePromise<ContentPublic> {
    return __request(OpenAPI, {
      method: "POST",
      url: "/api/v1/content/items",
      body: data.requestBody,
      mediaType: "application/json",
      errors: {
        422: "Validation Error",
      },
    })
  }

  public static updateContent(data: {
    contentId: number
    requestBody: ContentUpdateInput
  }): CancelablePromise<ContentPublic> {
    return __request(OpenAPI, {
      method: "PUT",
      url: "/api/v1/content/items/{content_id}",
      path: {
        content_id: data.contentId,
      },
      body: data.requestBody,
      mediaType: "application/json",
      errors: {
        422: "Validation Error",
      },
    })
  }

  public static deleteContent(data: {
    contentId: number
  }): CancelablePromise<MessagePublic> {
    return __request(OpenAPI, {
      method: "DELETE",
      url: "/api/v1/content/items/{content_id}",
      path: {
        content_id: data.contentId,
      },
      errors: {
        422: "Validation Error",
      },
    })
  }

  public static readTopics(data: {
    contentId: number
    skip?: number
    limit?: number
    behavior?: TopicBehavior
    search?: string
  }): CancelablePromise<TopicsPublic> {
    return __request(OpenAPI, {
      method: "GET",
      url: "/api/v1/content/items/{content_id}/topics",
      path: {
        content_id: data.contentId,
      },
      query: {
        skip: data.skip,
        limit: data.limit,
        behavior: data.behavior,
        search: data.search,
      },
      errors: {
        422: "Validation Error",
      },
    })
  }

  public static createTopic(data: {
    contentId: number
    requestBody: TopicCreateInput
  }): CancelablePromise<TopicPublic> {
    return __request(OpenAPI, {
      method: "POST",
      url: "/api/v1/content/items/{content_id}/topics",
      path: {
        content_id: data.contentId,
      },
      body: data.requestBody,
      mediaType: "application/json",
      errors: {
        422: "Validation Error",
      },
    })
  }

  public static updateTopic(data: {
    topicId: string
    requestBody: TopicUpdateInput
  }): CancelablePromise<TopicPublic> {
    return __request(OpenAPI, {
      method: "PUT",
      url: "/api/v1/content/topics/{topic_id}",
      path: {
        topic_id: data.topicId,
      },
      body: data.requestBody,
      mediaType: "application/json",
      errors: {
        422: "Validation Error",
      },
    })
  }

  public static deleteTopic(data: {
    topicId: string
  }): CancelablePromise<MessagePublic> {
    return __request(OpenAPI, {
      method: "DELETE",
      url: "/api/v1/content/topics/{topic_id}",
      path: {
        topic_id: data.topicId,
      },
      errors: {
        422: "Validation Error",
      },
    })
  }
}
