const statusLabels = {
    pending: "در انتظار بررسی",
    processing: "در حال پردازش",
    shipped: "ارسال شده",
    completed: "تکمیل شده",
    cancelled: "لغو شده"
};

const validTransitions = {
    pending: ["processing", "cancelled"],
    processing: ["shipped", "cancelled"],
    shipped: ["completed"],
    completed: [],
    cancelled: ["pending", "processing"],
}

let orders = [];

// Loads orders for the list in admin.html when the page opens.
async function loadOrders() {
    const response = await fetch("/api/orders");

    if (!response.ok) {
        document.getElementById("orders").textContent =
            "خطا در دریافت سفارش‌ها";
        return;
    }

    orders = await response.json();
    renderOrders();
}

// Builds the order cards and applies the status filter from admin.html.
function renderOrders() {
    const container = document.getElementById("orders");
    const selectedStatus = document.getElementById("status-filter").value;
    container.replaceChildren();
    container.classList.remove("empty-state");

    if (orders.length === 0) {
        container.textContent = "هنوز سفارشی ثبت نشده است.";
        return;
    }

    const visibleOrders = selectedStatus === "all"
        ? orders
        : orders.filter((order) => order.Status === selectedStatus);

    if (visibleOrders.length === 0) {
        container.textContent = "سفارشی با این وضعیت وجود ندارد.";
        container.classList.add("empty-state");
        return;
    }

    container.classList.remove("empty-state");

    for (const order of visibleOrders) {
        const orderElement = document.createElement("article");

        orderElement.className = "order-card";

        orderElement.innerHTML = `
            <h2>سفارش #${order.ID}</h2>

            <div class="order-meta">
                <div class="order-date">
                    <span class="order-date-icon" aria-hidden="true">▦</span>
                    <div>
                        <span class="order-meta-label">تاریخ سفارش</span>
                        <time datetime="${order.CreatedAt}">
                            ${new Date(order.CreatedAt).toLocaleDateString("fa-IR")}
                        </time>
                    </div>
                </div>

                <div class="order-status-field" data-status="${order.Status}">
                    <label for="order-status-${order.ID}">وضعیت سفارش</label>
                    <div class="order-status-control">
                        <svg
                            class="status-icon"
                            viewBox="0 0 24 24"
                            aria-hidden="true"
                            focusable="false"
                        >
                            <use href="#status-icon-${order.Status}"></use>
                        </svg>
                        <select
                            id="order-status-${order.ID}"
                            class="order-status"
                            data-order-id="${order.ID}"
                            data-current-status="${order.Status}"
                        >

                            ${[
                                order.Status,
                                ...validTransitions[order.Status]
                            ].map(status => `
                                <option
                                    value="${status}"
                                    ${order.Status === status ? "selected" : ""}
                                >
                                    ${statusLabels[status]}
                                </option>
                            `).join("")}
                        </select>
                    </div>
                </div>
            </div>

            <p>
                <strong>نام:</strong>
                ${order.CustomerName}
            </p>

            <p>
                <strong>تلفن:</strong>
                ${order.CustomerPhone}
            </p>

            <p>
                <strong>آدرس:</strong>
                ${order.CustomerAddress}
            </p>

            <h3>محصولات</h3>

            <div class="order-items"></div>

            <p>
                <strong>مبلغ کل:</strong>
                ${order.TotalPrice.toLocaleString()} تومان
            </p>
        `;

        const itemsContainer = orderElement.querySelector(".order-items");

        for (const item of order.Items) {
            const itemElement = document.createElement("div");

            itemElement.className = "order-item";

            const subtotal = item.UnitPrice * item.Quantity;

            itemElement.innerHTML = `
                <img
                    src="${item.Product.image_url}"
                    alt="${item.Product.name}"
                >

                <div>
                    <strong>${item.Product.name}</strong>

                    <p>
                        تعداد: ${item.Quantity}
                    </p>

                    <p>
                        قیمت واحد: ${item.UnitPrice.toLocaleString()} تومان
                    </p>

                    <p>
                        جمع: ${subtotal.toLocaleString()} تومان
                    </p>
                </div>
            `;

            itemsContainer.appendChild(itemElement);
        }

        container.appendChild(orderElement);
    }
}

// Shows the result message after an order status change.
function showNotification(message, type) {
    document.querySelectorAll(".notification").forEach((item) => item.remove());

    const notification = document.createElement("div");
    notification.className = `notification notification-${type}`;
    notification.textContent = message;
    notification.setAttribute("role", type === "error" ? "alert" : "status");

    document.body.appendChild(notification);

    setTimeout(() => {
        notification.remove();
    }, 3000);
}

// Handles changes to the status dropdowns rendered inside #orders.
document.addEventListener("change", async (event) => {
    if (!event.target.classList.contains("order-status")) {
        return;
    }

    const select = event.target;

    const orderID = select.dataset.orderId;
    const previousStatus = select.dataset.currentStatus;
    const newStatus = select.value;

    let response;
    select.disabled = true;
    try {
        response = await fetch(`/api/orders/${orderID}/status`, {
            method: "PATCH",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify({ status: newStatus })
        });
        if (!response.ok) {
            select.value = previousStatus;

            if (response.status === 409) {
                const message = await response.text();

                if (message === "Not enough stock\n") {
                    showNotification("موجودی کافی نیست", "error");
                } else if (message === "Invalid order status transition\n") {
                    showNotification("این تغییر وضعیت مجاز نیست", "error");
                } else {
                    showNotification("تغییر وضعیت سفارش ناموفق بود", "error");
                }
            }

            return;
        }

    } catch (error) {
        select.value = previousStatus;
        showNotification("ارتباط با سرور برقرار نشد", "error");
        console.error("Status update error:", error)
        return;

    } finally {
        select.disabled = false;
    }

    select.dataset.currentStatus = newStatus;
    const updatedOrder = orders.find(
        (order) => String(order.ID) === orderID
    );
    if (updatedOrder) {
        updatedOrder.Status = newStatus;
    }
    showNotification("وضعیت سفارش با موفقیت تغییر کرد", "success");
    renderOrders();
});

// Re-renders the order list when the filter select in admin.html changes.
document.getElementById("status-filter").addEventListener("change", renderOrders);
document.addEventListener("DOMContentLoaded", loadOrders);
