from django.contrib.auth.decorators import login_required
from django.core.paginator import Paginator
from django.db.models import Q
from django.http import JsonResponse
from django.shortcuts import render, redirect, get_object_or_404
from django.contrib.auth import login, authenticate, logout
from django.contrib.auth.forms import AuthenticationForm
from django.contrib.auth.views import LoginView as BaseLoginView

from . import models, forms


class LoginView(BaseLoginView):
    template_name = "login.html"
    authentication_form = AuthenticationForm


def signup_view(request):
    if request.method == "POST":
        form = forms.SignUpForm(request.POST)
        if form.is_valid():
            form.save()
            username = form.cleaned_data.get("username")
            raw_password = form.cleaned_data.get("password1")
            user = authenticate(username=username, password=raw_password)
            login(request, user)
            return redirect("home")
    else:
        form = forms.SignUpForm()
    return render(request, "signup.html", {"form": form})


@login_required
def home_view(request):
    query = request.GET.get("q")
    if query:
        announcements = models.Announcement.objects.filter(
            Q(book__title__icontains=query) | Q(book__author__icontains=query),
            status="open",
        ).select_related("book")
    else:
        announcements = models.Announcement.objects.filter(
            status="open"
        ).select_related("book")

    paginator = Paginator(announcements, 10)
    page_number = request.GET.get("page")
    page_obj = paginator.get_page(page_number)

    return render(request, "home.html", {"page_obj": page_obj, "query": query})


def logout_view(request):
    logout(request)
    return redirect("login")


@login_required
def manage_book_view(request, pk=None):
    if pk:
        book = get_object_or_404(models.Book, pk=pk, owner=request.user)
    else:
        book = None

    if request.method == "POST":
        form = forms.BookForm(request.POST, instance=book)

        if form.is_valid():
            book = form.save(commit=False)
            book.owner = request.user
            book.save()

            for img in request.FILES.getlist("images"):
                models.ImageBook.objects.create(book=book, image=img)

            return redirect("my_books")
    if request.method == "DELETE":
        book.images.all().delete()
        book.delete()
        return redirect("my_books")
    else:
        form = forms.BookForm(instance=book)
    images = book.images.all() if book and book.images.all().exists() else None
    return render(
        request, "manage_book.html", {"form": form, "book": book, "images": images}
    )


@login_required
def list_books_view(request):
    books = models.Book.objects.filter(owner=request.user)
    paginator = Paginator(books, 10)
    page_number = request.GET.get("page")
    page_obj = paginator.get_page(page_number)

    return render(request, "my_books.html", {"page_obj": page_obj})


@login_required
def delete_image(request, pk):
    image = get_object_or_404(models.ImageBook, pk=pk, book__owner=request.user)
    image.delete()
    return redirect("edit_book", pk=image.book.pk)


@login_required
def list_announcements_view(request):
    announcements = models.Announcement.objects.filter(
        created_by=request.user, status="open"
    ).select_related("book")

    paginator = Paginator(announcements, 10)
    page_number = request.GET.get("page")
    page_obj = paginator.get_page(page_number)
    return render(
        request, "list_announcements.html", {"page_obj": page_obj}
    )


@login_required
def manage_announcement_view(request, pk=None):
    if pk:
        announcement = get_object_or_404(
            models.Announcement, pk=pk, created_by=request.user
        )
        form = forms.AnnouncementForm(
            instance=announcement, initial={"user": request.user}
        )
    else:
        announcement = None
        form = forms.AnnouncementForm(initial={"user": request.user})

    if request.method == "POST":
        form = forms.AnnouncementForm(
            request.POST, instance=announcement, initial={"user": request.user}
        )
        if form.is_valid():
            announcement = form.save(commit=False)
            announcement.created_by = request.user
            announcement.save()
            return redirect("my_announcements")

    return render(
        request,
        "manage_announcements.html",
        {"form": form, "announcement": announcement},
    )


@login_required
def delete_announcement_view(request, pk):
    announcement = get_object_or_404(
        models.Announcement, pk=pk, created_by=request.user
    )
    if request.method == "POST":
        announcement.delete()
        return JsonResponse({"success": True})
    return JsonResponse(
        {"success": False, "error": "Request method not allowed"}, status=405
    )

@login_required
def detail_announcement_view(request, pk):
    announcement = get_object_or_404(models.Announcement, pk=pk)

    user_offer = models.OfferAnnouncement.objects.filter(announcement=announcement, user=request.user).first()

    if not user_offer and request.method == "POST":
        form = forms.OfferAnnouncementForm(request.POST, user=request.user, announcement=announcement)
        if form.is_valid():
            offer = form.save(commit=False)
            offer.user = request.user
            offer.announcement = announcement
            offer.save()
            return redirect("detail_announcement", pk=pk)
    else:
        form = forms.OfferAnnouncementForm(user=request.user, announcement=announcement) if not user_offer else None

    return render(request, "detail_announcement.html", {
        "announcement": announcement,
        "user_offer": user_offer,
        "form": form,
        "images": announcement.book.images.all()
    })