-- Таблица продуктов
CREATE TABLE IF NOT EXISTS products (
    product_id UUID PRIMARY KEY,
    product_name VARCHAR(255) NOT NULL,
    price BIGINT NOT NULL CHECK (price >= 0),
    rating SMALLINT CHECK (rating >= 0 AND rating <= 5),
    image TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Таблица корзины
CREATE TABLE IF NOT EXISTS cart (
    id UUID PRIMARY KEY,
    user_id VARCHAR(255) NOT NULL,
    product_id UUID NOT NULL,
    quantity INTEGER NOT NULL DEFAULT 1 CHECK (quantity > 0),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE,
    FOREIGN KEY (product_id) REFERENCES products(product_id) ON DELETE CASCADE,
    UNIQUE(user_id, product_id)
);

-- Таблица заказов (для будущих функций BuyFromCart и InstantBuy)
CREATE TABLE IF NOT EXISTS orders (
    order_id UUID PRIMARY KEY,
    user_id VARCHAR(255) NOT NULL,
    total_price BIGINT NOT NULL CHECK (total_price >= 0),
    ordered_at TIMESTAMP NOT NULL DEFAULT NOW(),
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
);

-- Таблица элементов заказа
CREATE TABLE IF NOT EXISTS order_items (
    id UUID PRIMARY KEY,
    order_id UUID NOT NULL,
    product_id UUID NOT NULL,
    quantity INTEGER NOT NULL CHECK (quantity > 0),
    price BIGINT NOT NULL CHECK (price >= 0),
    FOREIGN KEY (order_id) REFERENCES orders(order_id) ON DELETE CASCADE,
    FOREIGN KEY (product_id) REFERENCES products(product_id) ON DELETE RESTRICT
);

-- Индексы для оптимизации запросов
CREATE INDEX IF NOT EXISTS idx_cart_user_id ON cart(user_id);
CREATE INDEX IF NOT EXISTS idx_cart_product_id ON cart(product_id);
CREATE INDEX IF NOT EXISTS idx_orders_user_id ON orders(user_id);
CREATE INDEX IF NOT EXISTS idx_order_items_order_id ON order_items(order_id);

-- Добавим несколько тестовых продуктов
INSERT INTO products (product_id, product_name, price, rating, image) VALUES
    ('550e8400-e29b-41d4-a716-446655440001', 'Laptop Dell XPS 15', 150000, 5, 'https://example.com/laptop.jpg'),
    ('550e8400-e29b-41d4-a716-446655440002', 'iPhone 15 Pro', 120000, 5, 'https://example.com/iphone.jpg'),
    ('550e8400-e29b-41d4-a716-446655440003', 'Sony Headphones WH-1000XM5', 35000, 4, 'https://example.com/headphones.jpg'),
    ('550e8400-e29b-41d4-a716-446655440004', 'Samsung Galaxy Tab S9', 75000, 4, 'https://example.com/tablet.jpg'),
    ('550e8400-e29b-41d4-a716-446655440005', 'Apple Watch Series 9', 45000, 5, 'https://example.com/watch.jpg')
ON CONFLICT (product_id) DO NOTHING;
