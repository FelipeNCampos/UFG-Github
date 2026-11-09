from django.db import models
from django.contrib.auth.models import User


class ImageBook(models.Model):
    book = models.ForeignKey(
        "Book", on_delete=models.CASCADE, related_name="images", verbose_name="Livro"
    )
    image = models.ImageField(upload_to="books/", verbose_name="Imagem")
    created_at = models.DateTimeField(auto_now_add=True, verbose_name="Data de Criação")
    updated_at = models.DateTimeField(auto_now=True, verbose_name="Última Atualização")

    def __str__(self):
        return self.book.title

    class Meta:
        verbose_name = "Imagem do Livro"
        verbose_name_plural = "Imagens dos Livros"


class Book(models.Model):
    title = models.CharField(max_length=100, verbose_name="Título")
    author = models.CharField(max_length=100, verbose_name="Autor")
    edition = models.CharField(max_length=100, blank=True, verbose_name="Edição")
    year = models.IntegerField(blank=True, null=True, verbose_name="Ano")
    pages = models.IntegerField(blank=True, null=True, verbose_name="Páginas")
    publisher = models.CharField(max_length=100, verbose_name="Editora", blank=True)
    created_at = models.DateTimeField(auto_now_add=True, verbose_name="Data de Criação")
    updated_at = models.DateTimeField(auto_now=True, verbose_name="Última Atualização")
    owner = models.ForeignKey(
        User, on_delete=models.CASCADE, verbose_name="Proprietário"
    )

    def __str__(self):
        return self.title

    def get_first_image(self):
        return self.images.first()

    class Meta:
        verbose_name = "Livro"
        verbose_name_plural = "Livros"


class Announcement(models.Model):
    STATUS_CHOICES = (
        ("open", "Em Aberto"),
        ("finished", "Finalizado"),
    )
    title = models.CharField(max_length=255, verbose_name="Título", blank=True)
    book = models.ForeignKey(Book, on_delete=models.CASCADE, verbose_name="Livro")
    status = models.CharField(
        max_length=10, choices=STATUS_CHOICES, default="open", verbose_name="Status"
    )
    created_by = models.ForeignKey(
        User,
        on_delete=models.CASCADE,
        verbose_name="Criado por",
        related_name="announcements",
    )
    created_at = models.DateTimeField(auto_now_add=True, verbose_name="Data de Criação")
    updated_at = models.DateTimeField(auto_now=True, verbose_name="Última Atualização")

    def __str__(self):
        return f"Anúncio de {self.book.title}"

    class Meta:
        verbose_name = "Anúncio"
        verbose_name_plural = "Anúncios"


class Review(models.Model):
    book = models.ForeignKey(
        Book, on_delete=models.CASCADE, related_name="reviews", verbose_name="Livro"
    )
    user = models.ForeignKey(User, on_delete=models.CASCADE, verbose_name="Usuário")
    comment = models.TextField(blank=True, verbose_name="Comentário")
    rating = models.IntegerField(default=0, verbose_name="Nota")
    created_at = models.DateTimeField(auto_now_add=True, verbose_name="Data de Criação")

    def __str__(self):
        return f"{self.book.title} - {self.user.username}"

    class Meta:
        verbose_name = "Avaliação"
        verbose_name_plural = "Avaliações"


class OfferAnnouncement(models.Model):
    announcement = models.ForeignKey(
        'Announcement', on_delete=models.CASCADE, related_name='offers', verbose_name="Anúncio"
    )
    user = models.ForeignKey(
        User, on_delete=models.CASCADE, related_name='offers_made', verbose_name="Usuário"
    )
    offered_book = models.ForeignKey(
        'Book', on_delete=models.CASCADE, related_name='book_offers', verbose_name="Livro Ofertado"
    )
    comment = models.TextField(blank=True, verbose_name="Comentário")
    accepted = models.BooleanField(null=True, verbose_name="Oferta Aceita")
    created_at = models.DateTimeField(auto_now_add=True, verbose_name="Data de Criação")

    def __str__(self):
        return f"Oferta de {self.user.username} para {self.announcement.book.title}"

    class Meta:
        verbose_name = "Oferta de Anúncio"
        verbose_name_plural = "Ofertas de Anúncios"