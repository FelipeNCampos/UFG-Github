import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query"
import { Pencil, Plus, Trash2 } from "lucide-react"
import { useEffect, useMemo, useState } from "react"

import {
  ContentAdminService,
  type ActivityType,
  type ContentCreateInput,
  type ContentPublic,
  type TopicBehavior,
  type TopicCreateInput,
  type TopicPublic,
  type TopicUpdateInput,
} from "@/client/content-admin"
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from "@/components/ui/alert-dialog"
import { Button } from "@/components/ui/button"
import {
  Dialog,
  DialogClose,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog"
import { Input } from "@/components/ui/input"
import { LoadingButton } from "@/components/ui/loading-button"
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select"
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table"
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs"
import { Textarea } from "@/components/ui/textarea"
import useCustomToast from "@/hooks/useCustomToast"
import { handleError } from "@/utils"

type ContentFormState = {
  plan: string
  language: string
  day_id: string
  title: string
  content: string
}

type TopicFormState = {
  behavior: TopicBehavior
  activity_type: ActivityType
  order: string
  explanation_title: string
  explanation_content: string
  match_title: string
  match_description: string
  match_pairs: string
  fill_title: string
  fill_sentence: string
  fill_answer: string
  fill_options: string
  fill_explanation: string
  quizz_title: string
  quizz_question: string
  quizz_options: string
  prompt_title: string
  prompt_description: string
  prompt_ai_tool_context: string
  prompt_success_criteria: string
}

const PLAN_OPTIONS = ["base", "freelancer", "hub"] as const
const LANGUAGE_OPTIONS = ["pt-br", "en", "es", "fr"] as const

const defaultContentFormState: ContentFormState = {
  plan: "base",
  language: "pt-br",
  day_id: "1",
  title: "",
  content: "",
}

const defaultTopicFormState: TopicFormState = {
  behavior: "explicacao",
  activity_type: "MatchWords",
  order: "1",
  explanation_title: "",
  explanation_content: "",
  match_title: "",
  match_description: "",
  match_pairs: "",
  fill_title: "",
  fill_sentence: "",
  fill_answer: "",
  fill_options: "",
  fill_explanation: "",
  quizz_title: "",
  quizz_question: "",
  quizz_options: "",
  prompt_title: "",
  prompt_description: "",
  prompt_ai_tool_context: "",
  prompt_success_criteria: "",
}

function parseLines(value: string): string[] {
  return value
    .split("\n")
    .map((line) => line.trim())
    .filter(Boolean)
}

function buildTopicPayload(form: TopicFormState): {
  behavior: TopicBehavior
  activity_type?: ActivityType
  payload: Record<string, unknown>
} | null {
  if (form.behavior === "explicacao") {
    if (!form.explanation_title.trim() || !form.explanation_content.trim()) return null
    return {
      behavior: "explicacao",
      payload: {
        title: form.explanation_title.trim(),
        content: form.explanation_content.trim(),
      },
    }
  }

  if (form.activity_type === "MatchWords") {
    const pairs = parseLines(form.match_pairs)
      .map((line) => {
        const [text, ...matchParts] = line.split("|")
        return { text: text?.trim() ?? "", match: matchParts.join("|").trim() }
      })
      .filter((pair) => pair.text && pair.match)
    if (!form.match_title.trim() || !form.match_description.trim() || pairs.length === 0) return null
    return {
      behavior: "atividade",
      activity_type: "MatchWords",
      payload: {
        title: form.match_title.trim(),
        description: form.match_description.trim(),
        pairs,
      },
    }
  }

  if (form.activity_type === "FillBlanks") {
    const options = parseLines(form.fill_options)
    if (!form.fill_title.trim() || !form.fill_sentence.trim() || !form.fill_answer.trim() || options.length === 0 || !form.fill_explanation.trim()) return null
    return {
      behavior: "atividade",
      activity_type: "FillBlanks",
      payload: {
        title: form.fill_title.trim(),
        sentence: form.fill_sentence.trim(),
        answer: form.fill_answer.trim(),
        options,
        explanation: form.fill_explanation.trim(),
      },
    }
  }

  if (form.activity_type === "Quizz") {
    const options = parseLines(form.quizz_options)
      .map((line) => {
        const [text, isCorrectRaw, ...explanationParts] = line.split("|")
        return {
          text: text?.trim() ?? "",
          is_correct: (isCorrectRaw ?? "").trim().toLowerCase() === "true",
          explanation: explanationParts.join("|").trim(),
        }
      })
      .filter((item) => item.text && item.explanation)
    if (!form.quizz_title.trim() || !form.quizz_question.trim() || options.length === 0) return null
    return {
      behavior: "atividade",
      activity_type: "Quizz",
      payload: { title: form.quizz_title.trim(), question: form.quizz_question.trim(), options },
    }
  }

  if (!form.prompt_title.trim() || !form.prompt_description.trim() || !form.prompt_ai_tool_context.trim() || !form.prompt_success_criteria.trim()) return null
  return {
    behavior: "atividade",
    activity_type: "PromptTreiner",
    payload: {
      title: form.prompt_title.trim(),
      description: form.prompt_description.trim(),
      AIToolContext: form.prompt_ai_tool_context.trim(),
      successCriteria: form.prompt_success_criteria.trim(),
    },
  }
}

function getTopicDisplayTitle(topic: TopicPublic): string {
  if (topic.title) return topic.title
  const payloadTitle = topic.payload?.title
  return typeof payloadTitle === "string" ? payloadTitle : "-"
}

function getTopicDisplaySummary(topic: TopicPublic): string {
  if (topic.body) return topic.body
  const payload = topic.payload
  if (topic.behavior === "explicacao") return typeof payload?.content === "string" ? payload.content : "-"
  if (topic.activity_type === "MatchWords") return typeof payload?.description === "string" ? payload.description : "-"
  if (topic.activity_type === "FillBlanks") return typeof payload?.sentence === "string" ? payload.sentence : "-"
  if (topic.activity_type === "Quizz") return typeof payload?.question === "string" ? payload.question : "-"
  if (topic.activity_type === "PromptTreiner") return typeof payload?.description === "string" ? payload.description : "-"
  return "-"
}

function topicToForm(topic: TopicPublic): TopicFormState {
  const payload = topic.payload ?? {}
  const base = { ...defaultTopicFormState, behavior: topic.behavior, order: String(topic.order) }

  if (topic.behavior === "explicacao") {
    return {
      ...base,
      explanation_title: typeof payload.title === "string" ? payload.title : (topic.title ?? ""),
      explanation_content: typeof payload.content === "string" ? payload.content : (topic.body ?? ""),
    }
  }

  const activityType = topic.activity_type ?? "MatchWords"
  if (activityType === "MatchWords") {
    const pairs = Array.isArray(payload.pairs)
      ? payload.pairs
          .map((pair) => {
            if (!pair || typeof pair !== "object") return ""
            const text = (pair as { text?: unknown }).text
            const match = (pair as { match?: unknown }).match
            if (typeof text !== "string" || typeof match !== "string") return ""
            return `${text}|${match}`
          })
          .filter(Boolean)
          .join("\n")
      : ""
    return {
      ...base,
      activity_type: "MatchWords",
      match_title: typeof payload.title === "string" ? payload.title : (topic.title ?? ""),
      match_description: typeof payload.description === "string" ? payload.description : (topic.body ?? ""),
      match_pairs: pairs,
    }
  }
  if (activityType === "FillBlanks") {
    return {
      ...base,
      activity_type: "FillBlanks",
      fill_title: typeof payload.title === "string" ? payload.title : (topic.title ?? ""),
      fill_sentence: typeof payload.sentence === "string" ? payload.sentence : (topic.body ?? ""),
      fill_answer: typeof payload.answer === "string" ? payload.answer : "",
      fill_options: Array.isArray(payload.options) ? payload.options.map((opt) => String(opt)).join("\n") : "",
      fill_explanation: typeof payload.explanation === "string" ? payload.explanation : "",
    }
  }

  if (activityType === "Quizz") {
    const options = Array.isArray(payload.options)
      ? payload.options
          .map((opt) => {
            if (!opt || typeof opt !== "object") return ""
            const text = (opt as { text?: unknown }).text
            const isCorrect = (opt as { is_correct?: unknown }).is_correct
            const explanation = (opt as { explanation?: unknown }).explanation
            if (typeof text !== "string" || typeof explanation !== "string") return ""
            return `${text}|${Boolean(isCorrect)}|${explanation}`
          })
          .filter(Boolean)
          .join("\n")
      : ""
    return {
      ...base,
      activity_type: "Quizz",
      quizz_title: typeof payload.title === "string" ? payload.title : (topic.title ?? ""),
      quizz_question: typeof payload.question === "string" ? payload.question : (topic.body ?? ""),
      quizz_options: options,
    }
  }

  return {
    ...base,
    activity_type: "PromptTreiner",
    prompt_title: typeof payload.title === "string" ? payload.title : (topic.title ?? ""),
    prompt_description: typeof payload.description === "string" ? payload.description : (topic.body ?? ""),
    prompt_ai_tool_context: typeof payload.AIToolContext === "string" ? payload.AIToolContext : "",
    prompt_success_criteria: typeof payload.successCriteria === "string" ? payload.successCriteria : "",
  }
}

export function ContentManagement() {
  const queryClient = useQueryClient()
  const { showSuccessToast, showErrorToast } = useCustomToast()

  const [contentSearch, setContentSearch] = useState("")
  const [contentPlanFilter, setContentPlanFilter] = useState("all")
  const [contentLanguageFilter, setContentLanguageFilter] = useState("all")
  const [activeContentId, setActiveContentId] = useState<number | null>(null)

  const [topicSearch, setTopicSearch] = useState("")
  const [topicBehaviorFilter, setTopicBehaviorFilter] = useState<TopicBehavior | "all">("all")

  const [isContentDialogOpen, setIsContentDialogOpen] = useState(false)
  const [editingContent, setEditingContent] = useState<ContentPublic | null>(null)
  const [contentForm, setContentForm] = useState<ContentFormState>(defaultContentFormState)

  const [isTopicDialogOpen, setIsTopicDialogOpen] = useState(false)
  const [editingTopic, setEditingTopic] = useState<TopicPublic | null>(null)
  const [topicForm, setTopicForm] = useState<TopicFormState>(defaultTopicFormState)

  const [contentToDelete, setContentToDelete] = useState<ContentPublic | null>(null)
  const [topicToDelete, setTopicToDelete] = useState<TopicPublic | null>(null)

  const { data: contentsResponse, isLoading: isLoadingContents, isError: isContentsError } = useQuery({
    queryKey: ["admin-content-items", contentSearch, contentPlanFilter, contentLanguageFilter],
    queryFn: () =>
      ContentAdminService.readContents({
        skip: 0,
        limit: 200,
        search: contentSearch || undefined,
        plan: contentPlanFilter === "all" ? undefined : contentPlanFilter,
        language: contentLanguageFilter === "all" ? undefined : contentLanguageFilter,
      }),
  })

  const contents = contentsResponse?.data ?? []

  useEffect(() => {
    if (!contents.length) {
      setActiveContentId(null)
      return
    }
    const hasActive = contents.some((content) => content.id === activeContentId)
    if (!hasActive) setActiveContentId(contents[0].id)
  }, [contents, activeContentId])

  const activeContent = useMemo(() => contents.find((content) => content.id === activeContentId), [contents, activeContentId])

  const { data: topicsResponse, isLoading: isLoadingTopics, isError: isTopicsError } = useQuery({
    enabled: activeContentId !== null,
    queryKey: ["admin-content-topics", activeContentId, topicSearch, topicBehaviorFilter],
    queryFn: () => {
      if (activeContentId === null) throw new Error("Content id is required")
      return ContentAdminService.readTopics({
        contentId: activeContentId,
        skip: 0,
        limit: 200,
        search: topicSearch || undefined,
        behavior: topicBehaviorFilter === "all" ? undefined : topicBehaviorFilter,
      })
    },
  })

  const topics = topicsResponse?.data ?? []

  const refreshContents = () => queryClient.invalidateQueries({ queryKey: ["admin-content-items"] })
  const refreshTopics = () => queryClient.invalidateQueries({ queryKey: ["admin-content-topics"] })

  const createContentMutation = useMutation({
    mutationFn: (payload: ContentCreateInput) => ContentAdminService.createContent({ requestBody: payload }),
    onSuccess: (createdContent) => {
      showSuccessToast("Conteudo criado com sucesso")
      setIsContentDialogOpen(false)
      setEditingContent(null)
      setContentForm(defaultContentFormState)
      setActiveContentId(createdContent.id)
    },
    onError: handleError.bind(showErrorToast),
    onSettled: refreshContents,
  })

  const updateContentMutation = useMutation({
    mutationFn: (payload: { contentId: number; requestBody: ContentCreateInput }) => ContentAdminService.updateContent(payload),
    onSuccess: () => {
      showSuccessToast("Conteudo atualizado com sucesso")
      setIsContentDialogOpen(false)
      setEditingContent(null)
      setContentForm(defaultContentFormState)
    },
    onError: handleError.bind(showErrorToast),
    onSettled: refreshContents,
  })
  const deleteContentMutation = useMutation({
    mutationFn: (contentId: number) => ContentAdminService.deleteContent({ contentId }),
    onSuccess: () => {
      showSuccessToast("Conteudo removido com sucesso")
      setContentToDelete(null)
      setTopicSearch("")
      setTopicBehaviorFilter("all")
    },
    onError: handleError.bind(showErrorToast),
    onSettled: () => {
      refreshContents()
      refreshTopics()
    },
  })

  const createTopicMutation = useMutation({
    mutationFn: (payload: { contentId: number; requestBody: TopicCreateInput }) => ContentAdminService.createTopic(payload),
    onSuccess: () => {
      showSuccessToast("Topico criado com sucesso")
      setIsTopicDialogOpen(false)
      setEditingTopic(null)
      setTopicForm(defaultTopicFormState)
    },
    onError: handleError.bind(showErrorToast),
    onSettled: refreshTopics,
  })

  const updateTopicMutation = useMutation({
    mutationFn: (payload: { topicId: string; requestBody: TopicUpdateInput }) => ContentAdminService.updateTopic(payload),
    onSuccess: () => {
      showSuccessToast("Topico atualizado com sucesso")
      setIsTopicDialogOpen(false)
      setEditingTopic(null)
      setTopicForm(defaultTopicFormState)
    },
    onError: handleError.bind(showErrorToast),
    onSettled: refreshTopics,
  })

  const deleteTopicMutation = useMutation({
    mutationFn: (topicId: string) => ContentAdminService.deleteTopic({ topicId }),
    onSuccess: () => {
      showSuccessToast("Topico removido com sucesso")
      setTopicToDelete(null)
    },
    onError: handleError.bind(showErrorToast),
    onSettled: refreshTopics,
  })

  const openCreateContentDialog = () => {
    setEditingContent(null)
    setContentForm(defaultContentFormState)
    setIsContentDialogOpen(true)
  }

  const openEditContentDialog = (content: ContentPublic) => {
    setEditingContent(content)
    setContentForm({
      plan: content.plan,
      language: content.language,
      day_id: String(content.day_id),
      title: content.title,
      content: content.content,
    })
    setIsContentDialogOpen(true)
  }

  const openCreateTopicDialog = () => {
    setEditingTopic(null)
    setTopicForm(defaultTopicFormState)
    setIsTopicDialogOpen(true)
  }

  const openEditTopicDialog = (topic: TopicPublic) => {
    setEditingTopic(topic)
    setTopicForm(topicToForm(topic))
    setIsTopicDialogOpen(true)
  }

  const submitContent = () => {
    const dayId = Number.parseInt(contentForm.day_id, 10)
    if (!dayId || dayId < 1 || !contentForm.title.trim() || !contentForm.content.trim()) {
      showErrorToast("Preencha os campos obrigatorios do conteudo")
      return
    }
    const payload: ContentCreateInput = {
      plan: contentForm.plan.trim(),
      language: contentForm.language.trim(),
      day_id: dayId,
      title: contentForm.title.trim(),
      content: contentForm.content.trim(),
    }
    if (editingContent) {
      updateContentMutation.mutate({ contentId: editingContent.id, requestBody: payload })
      return
    }
    createContentMutation.mutate(payload)
  }

  const submitTopic = () => {
    if (activeContentId === null) {
      showErrorToast("Selecione um conteudo antes de criar topico")
      return
    }
    const order = Number.parseInt(topicForm.order, 10)
    if (!order || order < 1) {
      showErrorToast("Informe uma ordem valida")
      return
    }

    const built = buildTopicPayload(topicForm)
    if (!built) {
      showErrorToast("Preencha corretamente os campos do tipo de topico")
      return
    }

    if (editingTopic) {
      updateTopicMutation.mutate({
        topicId: editingTopic.id,
        requestBody: {
          behavior: built.behavior,
          activity_type: built.activity_type,
          payload: built.payload,
          order,
        },
      })
      return
    }

    createTopicMutation.mutate({
      contentId: activeContentId,
      requestBody: {
        content_id: activeContentId,
        behavior: built.behavior,
        activity_type: built.activity_type,
        payload: built.payload,
        order,
      },
    })
  }

  return (
    <div className="flex flex-col gap-6">
      <div className="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
        <div>
          <h1 className="text-2xl font-bold tracking-tight">Conteudo</h1>
          <p className="text-muted-foreground">Gerencie os conteudos e topicos por plano, idioma e dia.</p>
        </div>
        <Button className="w-full sm:w-auto" onClick={openCreateContentDialog}>
          <Plus className="mr-2 h-4 w-4" />Novo Conteudo
        </Button>
      </div>

      <div className="grid grid-cols-1 gap-3 md:grid-cols-3">
        <Input placeholder="Pesquisar conteudo..." value={contentSearch} onChange={(event) => setContentSearch(event.target.value)} />
        <Select value={contentPlanFilter} onValueChange={setContentPlanFilter}>
          <SelectTrigger><SelectValue placeholder="Plano" /></SelectTrigger>
          <SelectContent>
            <SelectItem value="all">Todos os planos</SelectItem>
            <SelectItem value="base">Base</SelectItem>
            <SelectItem value="freelancer">Freelancer</SelectItem>
            <SelectItem value="hub">Hub</SelectItem>
          </SelectContent>
        </Select>
        <Select value={contentLanguageFilter} onValueChange={setContentLanguageFilter}>
          <SelectTrigger><SelectValue placeholder="Idioma" /></SelectTrigger>
          <SelectContent>
            <SelectItem value="all">Todos os idiomas</SelectItem>
            {LANGUAGE_OPTIONS.map((language) => (
              <SelectItem key={language} value={language}>{language}</SelectItem>
            ))}
          </SelectContent>
        </Select>
      </div>

      {isLoadingContents ? <div className="rounded-md border p-6 text-sm text-muted-foreground">Carregando conteudos...</div> : null}
      {isContentsError ? <div className="rounded-md border border-destructive/30 bg-destructive/10 p-4 text-sm text-destructive">Falha ao carregar conteudos.</div> : null}
      {!isLoadingContents && !isContentsError && contents.length === 0 ? <div className="rounded-md border p-6 text-sm text-muted-foreground">Nenhum conteudo encontrado com os filtros atuais.</div> : null}

      {!isLoadingContents && !isContentsError && contents.length > 0 && activeContent ? (
      <Tabs
        value={activeContentId === null ? "" : String(activeContentId)}
        onValueChange={(value) => setActiveContentId(Number.parseInt(value, 10))}
      >
        <TabsList className="w-full justify-start overflow-x-auto">
          {contents.map((content) => (
              <TabsTrigger key={content.id} value={String(content.id)}>{content.plan}/{content.language} - Dia {content.day_id}</TabsTrigger>
          ))}
        </TabsList>

          <TabsContent value={String(activeContentId)} className="pt-4">
            <div className="flex flex-col gap-4 rounded-md border p-4">
              <div className="flex flex-col gap-3 sm:flex-row sm:items-start sm:justify-between">
                <div className="space-y-1">
                  <h3 className="text-lg font-semibold">{activeContent.title}</h3>
                  <p className="text-sm text-muted-foreground">Plano: {activeContent.plan} | Idioma: {activeContent.language} | Dia: {activeContent.day_id}</p>
                </div>
                <div className="flex gap-2">
                  <Button variant="outline" size="sm" onClick={() => openEditContentDialog(activeContent)}><Pencil className="mr-2 h-4 w-4" />Editar</Button>
                  <Button variant="destructive" size="sm" onClick={() => setContentToDelete(activeContent)}><Trash2 className="mr-2 h-4 w-4" />Remover</Button>
                </div>
              </div>
              <p className="whitespace-pre-wrap text-sm text-muted-foreground">{activeContent.content}</p>
            </div>

            <div className="mt-6 flex flex-col gap-4">
              <div className="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
                <div>
                  <h3 className="text-lg font-semibold">Topicos do conteudo</h3>
                  <p className="text-sm text-muted-foreground">Topicos tipados por comportamento e tipo de atividade.</p>
                </div>
                <Button className="w-full sm:w-auto" onClick={openCreateTopicDialog}><Plus className="mr-2 h-4 w-4" />Novo Topico</Button>
              </div>

              <div className="grid grid-cols-1 gap-3 md:grid-cols-2">
                <Input placeholder="Pesquisar topico..." value={topicSearch} onChange={(event) => setTopicSearch(event.target.value)} />
                <Select value={topicBehaviorFilter} onValueChange={(value) => setTopicBehaviorFilter(value as TopicBehavior | "all")}>
                  <SelectTrigger><SelectValue placeholder="Comportamento" /></SelectTrigger>
                  <SelectContent>
                    <SelectItem value="all">Todos</SelectItem>
                    <SelectItem value="explicacao">explicacao</SelectItem>
                    <SelectItem value="atividade">atividade</SelectItem>
                  </SelectContent>
                </Select>
              </div>

              {isLoadingTopics ? <div className="rounded-md border p-4 text-sm text-muted-foreground">Carregando topicos...</div> : null}
              {isTopicsError ? <div className="rounded-md border border-destructive/30 bg-destructive/10 p-4 text-sm text-destructive">Falha ao carregar topicos.</div> : null}

              {!isLoadingTopics && !isTopicsError ? (
                <div className="rounded-md border">
                  <Table>
                    <TableHeader>
                      <TableRow>
                        <TableHead>Ordem</TableHead>
                        <TableHead>Titulo</TableHead>
                        <TableHead>Comportamento</TableHead>
                        <TableHead>Tipo</TableHead>
                        <TableHead>Resumo</TableHead>
                        <TableHead className="text-right">Acoes</TableHead>
                      </TableRow>
                    </TableHeader>
                    <TableBody>
                      {topics.length === 0 ? (
                        <TableRow><TableCell colSpan={6} className="h-20 text-center text-muted-foreground">Nenhum topico encontrado.</TableCell></TableRow>
                      ) : (
                        topics.map((topic) => (
                          <TableRow key={topic.id}>
                            <TableCell>{topic.order}</TableCell>
                            <TableCell className="font-medium">{getTopicDisplayTitle(topic)}</TableCell>
                            <TableCell>{topic.behavior}</TableCell>
                            <TableCell>{topic.activity_type ?? "-"}</TableCell>
                            <TableCell className="max-w-2xl truncate text-muted-foreground">{getTopicDisplaySummary(topic)}</TableCell>
                            <TableCell className="text-right">
                              <div className="flex justify-end gap-2">
                                <Button variant="outline" size="sm" onClick={() => openEditTopicDialog(topic)}><Pencil className="mr-2 h-4 w-4" />Editar</Button>
                                <Button variant="destructive" size="sm" onClick={() => setTopicToDelete(topic)}><Trash2 className="mr-2 h-4 w-4" />Remover</Button>
                              </div>
                            </TableCell>
                          </TableRow>
                        ))
                      )}
                    </TableBody>
                  </Table>
                </div>
              ) : null}
            </div>
          </TabsContent>
        </Tabs>
      ) : null}

      <Dialog open={isContentDialogOpen} onOpenChange={setIsContentDialogOpen}>
        <DialogContent className="sm:max-w-2xl">
          <DialogHeader>
            <DialogTitle>{editingContent ? "Editar Conteudo" : "Novo Conteudo"}</DialogTitle>
            <DialogDescription>Defina plano, idioma, dia e texto principal do conteudo.</DialogDescription>
          </DialogHeader>
          <div className="grid grid-cols-1 gap-3 md:grid-cols-3">
            <Select value={contentForm.plan} onValueChange={(value) => setContentForm((prev) => ({ ...prev, plan: value }))}>
              <SelectTrigger><SelectValue placeholder="Plano" /></SelectTrigger>
              <SelectContent>{PLAN_OPTIONS.map((plan) => <SelectItem key={plan} value={plan}>{plan}</SelectItem>)}</SelectContent>
            </Select>
            <Input placeholder="Idioma" value={contentForm.language} onChange={(event) => setContentForm((prev) => ({ ...prev, language: event.target.value }))} />
            <Input type="number" min={1} placeholder="Dia" value={contentForm.day_id} onChange={(event) => setContentForm((prev) => ({ ...prev, day_id: event.target.value }))} />
          </div>
          <Input placeholder="Titulo" value={contentForm.title} onChange={(event) => setContentForm((prev) => ({ ...prev, title: event.target.value }))} />
          <Textarea placeholder="Conteudo" className="min-h-40" value={contentForm.content} onChange={(event) => setContentForm((prev) => ({ ...prev, content: event.target.value }))} />
          <DialogFooter>
            <DialogClose asChild><Button variant="outline" disabled={createContentMutation.isPending || updateContentMutation.isPending}>Cancelar</Button></DialogClose>
            <LoadingButton onClick={submitContent} loading={createContentMutation.isPending || updateContentMutation.isPending}>Salvar</LoadingButton>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      <Dialog open={isTopicDialogOpen} onOpenChange={setIsTopicDialogOpen}>
        <DialogContent className="sm:max-w-3xl">
          <DialogHeader>
            <DialogTitle>{editingTopic ? "Editar Topico" : "Novo Topico"}</DialogTitle>
            <DialogDescription>Defina o comportamento e preencha o payload do tipo selecionado.</DialogDescription>
          </DialogHeader>

          <div className="grid grid-cols-1 gap-3 md:grid-cols-3">
            <Select value={topicForm.behavior} onValueChange={(value) => setTopicForm((prev) => ({ ...prev, behavior: value as TopicBehavior }))}>
              <SelectTrigger><SelectValue placeholder="Comportamento" /></SelectTrigger>
              <SelectContent><SelectItem value="explicacao">explicacao</SelectItem><SelectItem value="atividade">atividade</SelectItem></SelectContent>
            </Select>
            {topicForm.behavior === "atividade" ? (
              <Select value={topicForm.activity_type} onValueChange={(value) => setTopicForm((prev) => ({ ...prev, activity_type: value as ActivityType }))}>
                <SelectTrigger><SelectValue placeholder="Tipo" /></SelectTrigger>
                <SelectContent><SelectItem value="MatchWords">MatchWords</SelectItem><SelectItem value="FillBlanks">FillBlanks</SelectItem><SelectItem value="Quizz">Quizz</SelectItem><SelectItem value="PromptTreiner">PromptTreiner</SelectItem></SelectContent>
              </Select>
            ) : <div />}
            <Input type="number" min={1} placeholder="Ordem" value={topicForm.order} onChange={(event) => setTopicForm((prev) => ({ ...prev, order: event.target.value }))} />
          </div>

          {topicForm.behavior === "explicacao" ? <div className="grid grid-cols-1 gap-3"><Input placeholder="Titulo da explicacao" value={topicForm.explanation_title} onChange={(event) => setTopicForm((prev) => ({ ...prev, explanation_title: event.target.value }))} /><Textarea placeholder="Conteudo da explicacao" className="min-h-36" value={topicForm.explanation_content} onChange={(event) => setTopicForm((prev) => ({ ...prev, explanation_content: event.target.value }))} /></div> : null}
          {topicForm.behavior === "atividade" && topicForm.activity_type === "MatchWords" ? <div className="grid grid-cols-1 gap-3"><Input placeholder="Titulo" value={topicForm.match_title} onChange={(event) => setTopicForm((prev) => ({ ...prev, match_title: event.target.value }))} /><Textarea placeholder="Descricao" className="min-h-24" value={topicForm.match_description} onChange={(event) => setTopicForm((prev) => ({ ...prev, match_description: event.target.value }))} /><Textarea placeholder="Pares (uma linha por par: texto|match)" className="min-h-32" value={topicForm.match_pairs} onChange={(event) => setTopicForm((prev) => ({ ...prev, match_pairs: event.target.value }))} /></div> : null}
          {topicForm.behavior === "atividade" && topicForm.activity_type === "FillBlanks" ? <div className="grid grid-cols-1 gap-3"><Input placeholder="Titulo" value={topicForm.fill_title} onChange={(event) => setTopicForm((prev) => ({ ...prev, fill_title: event.target.value }))} /><Textarea placeholder="Sentence" className="min-h-24" value={topicForm.fill_sentence} onChange={(event) => setTopicForm((prev) => ({ ...prev, fill_sentence: event.target.value }))} /><Input placeholder="Answer" value={topicForm.fill_answer} onChange={(event) => setTopicForm((prev) => ({ ...prev, fill_answer: event.target.value }))} /><Textarea placeholder="Options (uma por linha)" className="min-h-24" value={topicForm.fill_options} onChange={(event) => setTopicForm((prev) => ({ ...prev, fill_options: event.target.value }))} /><Textarea placeholder="Explanation" className="min-h-24" value={topicForm.fill_explanation} onChange={(event) => setTopicForm((prev) => ({ ...prev, fill_explanation: event.target.value }))} /></div> : null}
          {topicForm.behavior === "atividade" && topicForm.activity_type === "Quizz" ? <div className="grid grid-cols-1 gap-3"><Input placeholder="Titulo" value={topicForm.quizz_title} onChange={(event) => setTopicForm((prev) => ({ ...prev, quizz_title: event.target.value }))} /><Textarea placeholder="Question" className="min-h-24" value={topicForm.quizz_question} onChange={(event) => setTopicForm((prev) => ({ ...prev, quizz_question: event.target.value }))} /><Textarea placeholder="Options (uma por linha: text|true|explanation)" className="min-h-32" value={topicForm.quizz_options} onChange={(event) => setTopicForm((prev) => ({ ...prev, quizz_options: event.target.value }))} /></div> : null}
          {topicForm.behavior === "atividade" && topicForm.activity_type === "PromptTreiner" ? <div className="grid grid-cols-1 gap-3"><Input placeholder="Titulo" value={topicForm.prompt_title} onChange={(event) => setTopicForm((prev) => ({ ...prev, prompt_title: event.target.value }))} /><Textarea placeholder="Description" className="min-h-24" value={topicForm.prompt_description} onChange={(event) => setTopicForm((prev) => ({ ...prev, prompt_description: event.target.value }))} /><Textarea placeholder="AIToolContext" className="min-h-24" value={topicForm.prompt_ai_tool_context} onChange={(event) => setTopicForm((prev) => ({ ...prev, prompt_ai_tool_context: event.target.value }))} /><Textarea placeholder="successCriteria" className="min-h-24" value={topicForm.prompt_success_criteria} onChange={(event) => setTopicForm((prev) => ({ ...prev, prompt_success_criteria: event.target.value }))} /></div> : null}

          <DialogFooter>
            <DialogClose asChild><Button variant="outline" disabled={createTopicMutation.isPending || updateTopicMutation.isPending}>Cancelar</Button></DialogClose>
            <LoadingButton onClick={submitTopic} loading={createTopicMutation.isPending || updateTopicMutation.isPending}>Salvar</LoadingButton>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      <AlertDialog open={Boolean(contentToDelete)} onOpenChange={(open) => { if (!open) setContentToDelete(null) }}>
        <AlertDialogContent><AlertDialogHeader><AlertDialogTitle>Remover Conteudo</AlertDialogTitle><AlertDialogDescription>Esta acao remove o conteudo e todos os topicos vinculados. Deseja continuar?</AlertDialogDescription></AlertDialogHeader><AlertDialogFooter><AlertDialogCancel disabled={deleteContentMutation.isPending}>Cancelar</AlertDialogCancel><AlertDialogAction onClick={() => contentToDelete && deleteContentMutation.mutate(contentToDelete.id)} className="bg-destructive hover:bg-destructive/90">{deleteContentMutation.isPending ? "Removendo..." : "Remover"}</AlertDialogAction></AlertDialogFooter></AlertDialogContent>
      </AlertDialog>

      <AlertDialog open={Boolean(topicToDelete)} onOpenChange={(open) => { if (!open) setTopicToDelete(null) }}>
        <AlertDialogContent><AlertDialogHeader><AlertDialogTitle>Remover Topico</AlertDialogTitle><AlertDialogDescription>Deseja remover este topico do conteudo?</AlertDialogDescription></AlertDialogHeader><AlertDialogFooter><AlertDialogCancel disabled={deleteTopicMutation.isPending}>Cancelar</AlertDialogCancel><AlertDialogAction onClick={() => topicToDelete && deleteTopicMutation.mutate(topicToDelete.id)} className="bg-destructive hover:bg-destructive/90">{deleteTopicMutation.isPending ? "Removendo..." : "Remover"}</AlertDialogAction></AlertDialogFooter></AlertDialogContent>
      </AlertDialog>
    </div>
  )
}

export default ContentManagement
