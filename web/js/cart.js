// Fills the cart area in cart.html, including its quantity buttons.
async function loadCart() {
    const cart = JSON.parse(localStorage.getItem("cart")) || [];
    const container = document.getElementById("cart");
    const checkoutForm = document.getElementById("checkout-form");
    container.innerHTML = "";
    checkoutForm.hidden = cart.length === 0;

    if (cart.length === 0) {
        container.innerHTML = `
            <div class="empty-state">
                <p>سبد خرید شما خالی است.</p>
                <a class="button" href="/">دیدن محصولات</a>
            </div>
        `;
        return;
    }

    let cartTotal = 0;
    for (const item of cart) {
        const response = await fetch(`/api/products/${item.id}`);
        if (!response.ok) {
            container.innerHTML = "<p>دریافت اطلاعات محصولات ناموفق بود.</p>";
            return;
        }

        const product = await response.json();
        const total = product.price * item.quantity;

        cartTotal += total;

        const element = document.createElement("div");
        element.className = "cart-item";
        element.innerHTML = `
            <h2>${product.name}</h2>
            <p>قیمت: ${product.price.toLocaleString()} تومان</p>
            <!-- These buttons call changeQuantity() in /js/cart.js. -->
            <div class="quantity">
                <button
                    class="button-secondary"
                    aria-label="افزایش تعداد"
                    onclick="changeQuantity(${item.id}, +1)"
                >
                    +
                </button>
                <span>${item.quantity}</span>
                <button
                    class="button-secondary"
                    aria-label="کاهش تعداد"
                    onclick="changeQuantity(${item.id}, -1)"
                >
                    −
                </button>
            </div>
            <p class="item-total">جمع: ${total.toLocaleString()} تومان</p>
        `;

        container.appendChild(element);
    }

    container.innerHTML += `
        <div class="cart-total">مجموع: ${cartTotal.toLocaleString()} تومان</div>
    `
}

// Handles clicks on the cart's plus and minus buttons.
function changeQuantity(id, amount) {
    const cart = JSON.parse(localStorage.getItem("cart")) || [];

    const item = cart.find(item => item.id === id);
    if (!item) return;

    item.quantity += amount;

    if (item.quantity <= 0) {
        const index = cart.indexOf(item)
        cart.splice(index, 1)
    }

    localStorage.setItem("cart", JSON.stringify(cart));
    loadCart();
}

// Shows checkout success or error messages on the cart page.
function showNotification(message, type) {
    document.querySelectorAll(".notification").forEach((item) => item.remove());

    const notification = document.createElement("div");
    notification.className = `notification notification-${type}`;
    notification.textContent = message;
    notification.setAttribute("role", type === "error" ? "alert" : "status");
    document.body.appendChild(notification);

    setTimeout(() => notification.remove(), 3000);
}

// Called when the checkout form in cart.html is submitted.
async function checkout(event) {
    event.preventDefault();

    const cart = JSON.parse(localStorage.getItem("cart")) || [];

    if (cart.length === 0) {
        showNotification("سبد خرید شما خالی است.", "error");
        return;
    }

    const customerName = document.getElementById("customer-name").value;
    const customerPhone = document.getElementById("customer-phone").value;
    const customerAddress = document.getElementById("customer-address").value;

    let response;
    try {
        response = await fetch("/api/orders", {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify({
                customer_name: customerName,
                customer_phone: customerPhone,
                customer_address: customerAddress,
                items: cart
            })
        });
    } catch (error) {
        showNotification("ارتباط با سرور برقرار نشد. دوباره تلاش کنید.", "error");
        console.error("Order submission error:", error);
        return;
    }

    if (!response.ok) {
        const message = await response.text();
        showNotification(
            message || "ثبت سفارش ناموفق بود. لطفاً دوباره تلاش کنید.",
            "error"
        );
        return;
    }

    localStorage.removeItem("cart");
    loadCart();
    showNotification("سفارش شما با موفقیت ثبت شد.", "success");
}

document.addEventListener("DOMContentLoaded", loadCart);
