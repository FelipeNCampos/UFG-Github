(function () {
  "use strict";

  /**
   * Caminhos relativos à pasta do index.html.
   * Vídeos: interior-dark-video.mp4 (hero), video-camadas.mp4 (processo).
   * Portfólio: sete imagens — renomeie os arquivos locais para bater com a lista
   * ou edite os nomes abaixo (.jpg, .jpeg, .png ou .webp).
   * Retrato "Sobre": defina aboutPortrait; se não existir arquivo, use uma das fotos do portfólio.
   * WhatsApp: substitua DDI DDD número em whatsappUrl (somente dígitos após wa.me/).
   */
  var ASSET_BASE = "src/assets/";
  var CONFIG = {
    heroVideo: ASSET_BASE + "interior-dark-video.mp4",
    layersVideo: ASSET_BASE + "video-camadas.mp4",
    portfolioImages: [
      ASSET_BASE + "portfolio-01.jpg",
      ASSET_BASE + "portfolio-02.jpg",
      ASSET_BASE + "portfolio-03.jpg",
      ASSET_BASE + "portfolio-04.jpg",
      ASSET_BASE + "portfolio-05.jpg",
      ASSET_BASE + "portfolio-06.jpg",
      ASSET_BASE + "portfolio-07.jpg",
    ],
    aboutPortrait: ASSET_BASE + "marina-portrait.jpg",
    whatsappUrl: "https://wa.me/5511998877665",
  };

  var portfolioProjects = [
    { label: "Residência Alphaville", title: "Oasis Urbano" },
    { label: "Cobertura Jardins", title: "Luz Contida" },
    { label: "Casa Itu", title: "Horizonte Fechado" },
    { label: "Escritório Faria Lima", title: "Linhas de Negócio" },
    { label: "Apartamento Pinheiros", title: "Galeria Habitável" },
    { label: "Loft Vila Madalena", title: "Brisa e Bronze" },
    { label: "Estúdio de Materiais", title: "Camadas e Textura" },
  ];

  function qs(sel, root) {
    return (root || document).querySelector(sel);
  }

  function qsa(sel, root) {
    return Array.prototype.slice.call((root || document).querySelectorAll(sel));
  }

  function initVideos() {
    var hero = qs("[data-hero-video]");
    var layers = qs("[data-layers-video]");
    if (hero) {
      hero.src = CONFIG.heroVideo;
      hero.play().catch(function () {});
    }
    if (layers) {
      layers.src = CONFIG.layersVideo;
      layers.play().catch(function () {});
    }
  }

  function initPortfolioImages() {
    var items = qsa("[data-portfolio-item]");
    items.forEach(function (el, i) {
      var img = el.querySelector("img");
      if (!img || !CONFIG.portfolioImages[i]) return;
      img.src = CONFIG.portfolioImages[i];
      var cap = portfolioProjects[i];
      if (cap) {
        var label = el.querySelector("[data-portfolio-label]");
        var title = el.querySelector("[data-portfolio-title]");
        if (label) label.textContent = cap.label;
        if (title) title.textContent = cap.title;
      }
    });
  }

  function initAboutImage() {
    var img = qs("[data-about-portrait]");
    if (img) img.src = CONFIG.aboutPortrait;
  }

  function initWhatsAppLinks() {
    qsa("[data-whatsapp]").forEach(function (a) {
      a.href = CONFIG.whatsappUrl;
      a.setAttribute("rel", "noopener noreferrer");
      a.setAttribute("target", "_blank");
    });
  }

  function initHeader() {
    var header = qs("[data-site-header]");
    var toggle = qs("[data-menu-toggle]");
    var mobile = qs("[data-nav-mobile]");
    var mobileLinks = mobile ? qsa("a", mobile) : [];

    function setOpen(open) {
      if (!toggle || !mobile) return;
      toggle.setAttribute("aria-expanded", open ? "true" : "false");
      mobile.classList.toggle("is-open", open);
      document.body.style.overflow = open ? "hidden" : "";
    }

    if (toggle && mobile) {
      toggle.addEventListener("click", function () {
        var open = toggle.getAttribute("aria-expanded") !== "true";
        setOpen(open);
      });
      mobileLinks.forEach(function (link) {
        link.addEventListener("click", function () {
          setOpen(false);
        });
      });
    }

    function onScroll() {
      if (!header) return;
      header.classList.toggle("is-scrolled", window.scrollY > 24);
    }

    onScroll();
    window.addEventListener("scroll", onScroll, { passive: true });
  }

  function initActiveNav() {
    var links = qsa("[data-nav-link]");
    if (!links.length) return;

    var sections = links
      .map(function (a) {
        var id = a.getAttribute("href");
        if (!id || id.charAt(0) !== "#") return null;
        var el = document.getElementById(id.slice(1));
        return el ? { id: id, el: el } : null;
      })
      .filter(Boolean);

    function update() {
      var y = window.scrollY + 120;
      var current = sections[0] && sections[0].id;
      for (var i = 0; i < sections.length; i++) {
        var sec = sections[i];
        var rect = sec.el.getBoundingClientRect();
        var top = rect.top + window.scrollY;
        if (top <= y) current = sec.id;
      }
      links.forEach(function (a) {
        a.removeAttribute("aria-current");
        if (a.getAttribute("href") === current) a.setAttribute("aria-current", "page");
      });
    }

    window.addEventListener("scroll", update, { passive: true });
    update();
  }

  function initReveal() {
    if (window.matchMedia("(prefers-reduced-motion: reduce)").matches) {
      qsa(".reveal").forEach(function (el) {
        el.classList.add("is-visible");
      });
      return;
    }
    var els = qsa(".reveal");
    if (!els.length || !("IntersectionObserver" in window)) {
      els.forEach(function (el) {
        el.classList.add("is-visible");
      });
      return;
    }
    var io = new IntersectionObserver(
      function (entries) {
        entries.forEach(function (entry) {
          if (entry.isIntersecting) {
            entry.target.classList.add("is-visible");
            io.unobserve(entry.target);
          }
        });
      },
      { rootMargin: "0px 0px -8% 0px", threshold: 0.08 }
    );
    els.forEach(function (el) {
      io.observe(el);
    });
  }

  function validateEmail(v) {
    return /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(String(v).trim());
  }

  function initForm() {
    var form = qs("[data-contact-form]");
    if (!form) return;

    var status = qs("[data-form-status]");
    var name = qs("#field-name");
    var email = qs("#field-email");
    var message = qs("#field-message");

    function setInvalid(field, messageText, invalid) {
      var wrap = field && field.closest(".form-field");
      if (!wrap) return;
      wrap.classList.toggle("is-invalid", !!invalid);
      var err = wrap.querySelector(".form-error");
      if (err) err.textContent = messageText || "";
    }

    form.addEventListener("submit", function (e) {
      e.preventDefault();
      var ok = true;
      setInvalid(name, "", false);
      setInvalid(email, "", false);
      setInvalid(message, "", false);
      if (status) {
        status.textContent = "";
        status.className = "form-status";
      }

      if (!name || !String(name.value).trim()) {
        setInvalid(name, "Informe seu nome.", true);
        ok = false;
      }
      if (!email || !validateEmail(email.value)) {
        setInvalid(email, "Informe um e-mail válido.", true);
        ok = false;
      }
      if (!message || !String(message.value).trim()) {
        setInvalid(message, "Escreva uma mensagem breve.", true);
        ok = false;
      }
      if (!ok) return;

      if (status) {
        status.textContent = "Mensagem preparada. Abrindo seu cliente de e-mail…";
        status.className = "form-status";
      }

      var interest = qs("#field-interest");
      var subject = encodeURIComponent("Contato — Studio Aura");
      var body = encodeURIComponent(
        "Nome: " +
          String(name.value).trim() +
          "\nE-mail: " +
          String(email.value).trim() +
          "\nInteresse: " +
          (interest ? interest.value : "") +
          "\n\n" +
          String(message.value).trim()
      );
      window.location.href = "mailto:contato@studioaura.com.br?subject=" + subject + "&body=" + body;

      if (status) {
        status.textContent =
          "Se nada abriu automaticamente, envie para contato@studioaura.com.br.";
        status.className = "form-status form-status--success";
      }
      form.reset();
    });
  }

  document.addEventListener("DOMContentLoaded", function () {
    initVideos();
    initPortfolioImages();
    initAboutImage();
    initWhatsAppLinks();
    initHeader();
    initActiveNav();
    initReveal();
    initForm();
  });
})();
