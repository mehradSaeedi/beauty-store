// Fills the product grid in index.html.
async function loadProducts() {
    const response = await fetch("api/products");
    const products = await response.json();
    const container = document.getElementById("products");

    products.forEach(product => {
        const element = document.createElement("div");
        element.className = "product-card";
        element.innerHTML = `
            <a href="/product.html?id=${product.ID}" aria-label="مشاهده ${product.name}">
                <div class="product-card-image">
                    <img src="${product.image_url}" alt="${product.name}">
                </div>
                <div class="product-card-copy">
                    <span class="product-kicker">LaLa Marnia · مراقبت روزانه</span>
                    <h2>${product.name}</h2>
                    <p>${product.description}</p>
                    <span class="product-card-bottom">
                        <span class="price">
                            ${product.price.toLocaleString()} تومان
                        </span>
                        <span class="product-arrow" aria-hidden="true">←</span>
                    </span>
                </div>
            </a>
        `;

        container.appendChild(element)
    })
}

loadProducts();
