// Loads the product detail and builds its button in product.html.
async function loadProduct() {
    const params = new URLSearchParams(window.location.search);
    const id = params.get("id");
    const response = await fetch(`/api/products/${id}`);
    const product = await response.json();
    const container = document.getElementById("product");
    container.innerHTML = `
        <div class="product-detail-image">
            <img src="${product.image_url}" alt="${product.name}">
        </div>
        <section>
            <span class="eyebrow">LaLa Marnia · انتخابی برای مراقبت تو</span>
            <h1>${product.name}</h1>
            <p class="product-description">${product.description}</p>
            <p class="price product-detail-price">
                ${product.price.toLocaleString()} <span>تومان</span>
            </p>
            <!-- This button calls addToCart() from /js/product.js. -->
            <button
                class="add-to-cart"
                onclick="addToCart(${product.ID})"
            >
                افزودن به سبد خرید <span aria-hidden="true">←</span>
            </button>
            <p class="product-assurance">
                <span aria-hidden="true">✓</span>
                سفارش آسان و خرید مطمئن
            </p>
        </section>
    `;
}

loadProduct();

// Shows the message after the product button is used.
function showNotification(message, type) {
    document.querySelectorAll(".notification").forEach((item) => item.remove());

    const notification = document.createElement("div");
    notification.className = `notification notification-${type}`;
    notification.textContent = message;
    notification.setAttribute("role", type === "error" ? "alert" : "status");
    document.body.appendChild(notification);

    setTimeout(() => notification.remove(), 3000);
}

// This is the add-to-cart button's click handler.
function addToCart(id) {
    let cart = JSON.parse(localStorage.getItem("cart")) || [];
    const existingItem = cart.find(item => item.id === id);

    if (existingItem) {
        existingItem.quantity++;
    } else {
        cart.push({
            id: id,
            quantity: 1
        });
    }

    localStorage.setItem("cart", JSON.stringify(cart));
    showNotification("محصول به سبد خرید اضافه شد.", "success");
}
