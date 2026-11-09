from django.shortcuts import render
from django.urls import path
from . import views
from django.contrib.auth.views import LogoutView

urlpatterns = [
    path('login/', views.LoginView.as_view(), name='login'),
    path('logout/', views.logout_view, name='logout'),
    path('signup/', views.signup_view, name='signup'),
    path('', views.home_view, name="home"),
    path('livros/', views.manage_book_view, name='create_book'),
    path('livros/edit/<int:pk>/', views.manage_book_view, name='edit_book'),
    path('livros/delete/<int:pk>/', views.manage_book_view, name='delete_book'),
    path('imagens/delete/<int:pk>/', views.delete_image, name='delete_image'),
    path("meus-anuncios/", views.list_announcements_view, name="my_announcements"),
    path("anuncios/", views.manage_announcement_view, name="manage_announcement"),
    path("detalhe-anuncio/<int:pk>", views.detail_announcement_view, name="detail_announcement"),
    path("anuncios/edit/<int:pk>/", views.manage_announcement_view, name="edit_announcement"),
    path("anuncios/delete/<int:pk>/", views.delete_announcement_view, name="delete_announcement"),
    path('meus-livros/', views.list_books_view, name='my_books'),
]
