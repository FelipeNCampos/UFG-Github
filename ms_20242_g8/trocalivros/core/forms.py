from django import forms
from django.contrib.auth.forms import (
    AuthenticationForm as BaseAuthenticationForm,
    UserCreationForm,
)
from django.contrib.auth import models as auth_models
from . import models


class AuthenticationForm(BaseAuthenticationForm):
    username = forms.CharField(
        widget=forms.TextInput(
            attrs={"class": "form-control", "placeholder": "Digite seu usuário"}
        )
    )
    password = forms.CharField(
        widget=forms.PasswordInput(
            attrs={"class": "form-control", "placeholder": "Digite sua senha"}
        )
    )


class SignUpForm(UserCreationForm):
    email = forms.EmailField(
        widget=forms.EmailInput(
            attrs={"class": "form-control", "placeholder": "Digite seu e-mail"}
        )
    )

    class Meta:
        model = auth_models.User
        fields = ("username", "email", "password1", "password2")


class BookForm(forms.ModelForm):
    class Meta:
        model = models.Book
        fields = ["title", "author", "edition", "year", "pages", "publisher"]


class AnnouncementForm(forms.ModelForm):
    title = forms.CharField(
        label="Título",
        widget=forms.TextInput(
            attrs={
                "class": "form-control",
                "placeholder": "Escreva algo como: Troco livro de matemática por livro de física",
            }
        ),
    )
    book = forms.ModelChoiceField(
        label="Livro",
        help_text="Selecione um livro para trocar",
        queryset=models.Book.objects.none(),
        widget=forms.Select(attrs={"class": "form-control"}),
    )

    class Meta:
        model = models.Announcement
        fields = ["title", "book"]

    def __init__(self, *args, **kwargs):
        super().__init__(*args, **kwargs)
        if "user" in self.initial:
            user = self.initial["user"]
            self.fields["book"].queryset = models.Book.objects.filter(
                owner=user
            ).exclude(
                pk__in=models.Announcement.objects.filter(
                    created_by=user, status="open"
                ).values_list("book_id", flat=True)
            )
            if self.instance and self.instance.pk:
                self.fields["book"].queryset |= models.Book.objects.filter(pk=self.instance.book.pk)


class OfferAnnouncementForm(forms.ModelForm):
    comment = forms.CharField(
        label="Comentário",
        widget=forms.Textarea(attrs={"class": "form-control", "placeholder": "Escreva um comentário para sua oferta"}),
        required=False
    )

    offered_book = forms.ModelChoiceField(
        queryset=models.Book.objects.none(),
        widget=forms.Select(attrs={"class": "form-control"}),
        label="Selecione um de seus livros"
    )

    class Meta:
        model = models.OfferAnnouncement
        fields = ["offered_book", "comment"]

    def __init__(self, *args, **kwargs):
        user = kwargs.pop('user', None)
        announcement = kwargs.pop('announcement', None)
        super().__init__(*args, **kwargs)

        if user:
            self.fields['offered_book'].queryset = models.Book.objects.filter(owner=user).exclude(pk=announcement.book.pk)