import { useEffect, useState } from "react";
import "./App.css";

const API = "http://localhost:9000/api";

function App() {
  const [materials, setMaterials] = useState([]);
  const [products, setProducts] = useState([]);
  const [workOrders, setWorkOrders] = useState([]);
  const [salesOrders, setSalesOrders] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  const [selectedProduct, setSelectedProduct] = useState("");
  const [productionQuantity, setProductionQuantity] = useState(1);
  const [bomPreview, setBomPreview] = useState(null);
  const [actionLoading, setActionLoading] = useState(false);

  const [salesCustomer, setSalesCustomer] = useState("");
  const [salesProduct, setSalesProduct] = useState("");
  const [salesQuantity, setSalesQuantity] = useState(1);

  // Mengambil data dari backend
  const loadData = async () => {
    try {
      setLoading(true);
      setError("");

      const [materialRes, productRes, orderRes, salesOrderRes] = await Promise.all([
        fetch(`${API}/inventory/materials`),
        fetch(`${API}/inventory/products`),
        fetch(`${API}/work-orders`),
        fetch(`${API}/sales-orders`),
      ]);

      if (!materialRes.ok || !productRes.ok || !orderRes.ok || !salesOrderRes.ok) {
        throw new Error("Gagal mengambil data dari backend.");
      }

      const [materialData, productData, orderData, salesOrderData] = await Promise.all([
        materialRes.json(),
        productRes.json(),
        orderRes.json(),
        salesOrderRes.json(),
      ]);

      setMaterials(materialData.data || []);
      setProducts(productData.data || []);
      setWorkOrders(orderData.data || []);
      setSalesOrders(salesOrderData.data || []);
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadData();
  }, []);

  // Preview kebutuhan bahan baku berdasarkan BOM
  const handleBOMPreview = async () => {
    if (!selectedProduct || Number(productionQuantity) <= 0) {
      alert("Pilih produk dan masukkan jumlah produksi yang valid!");
      return;
    }

    try {
      setActionLoading(true);
      setBomPreview(null);

      const response = await fetch(
        `${API}/products/${selectedProduct}/bom-explosion`,
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify({
            quantity: Number(productionQuantity),
          }),
        }
      );

      const result = await response.json();

      if (!response.ok) {
        throw new Error(result.message || "Gagal melakukan BOM Explosion");
      }

      setBomPreview(result.data);
    } catch (err) {
      alert(err.message);
    } finally {
      setActionLoading(false);
    }
  };

  // Membuat Work Order
  const handleCreateWorkOrder = async () => {
    if (!selectedProduct || Number(productionQuantity) <= 0) {
      alert("Pilih produk dan masukkan jumlah produksi yang valid!");
      return;
    }

    try {
      setActionLoading(true);

      const response = await fetch(`${API}/work-orders`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          product_id: Number(selectedProduct),
          quantity: Number(productionQuantity),
        }),
      });

      const result = await response.json();

      if (!response.ok) {
        throw new Error(result.message || "Gagal membuat Work Order");
      }

      alert("Work Order berhasil dibuat!");

      setBomPreview(null);
      setSelectedProduct("");
      setProductionQuantity(1);

      await loadData();
    } catch (err) {
      alert(err.message);
    } finally {
      setActionLoading(false);
    }
  };

  // Membuat Sales Order dan reservasi stok produk
  const handleCreateSalesOrder = async () => {
    if (!salesCustomer.trim() || !salesProduct || Number(salesQuantity) <= 0) {
      alert("Lengkapi nama pelanggan, produk, dan quantity yang valid!");
      return;
    }

    try {
      setActionLoading(true);
      const response = await fetch(`${API}/sales-orders`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          customer_name: salesCustomer.trim(),
          items: [{
            product_id: Number(salesProduct),
            quantity: Number(salesQuantity),
          }],
        }),
      });
      const result = await response.json();
      if (!response.ok) {
        throw new Error(result.error || result.message || "Gagal membuat Sales Order");
      }

      alert("Sales Order berhasil dibuat dan stok berhasil direservasi!");
      setSalesCustomer("");
      setSalesProduct("");
      setSalesQuantity(1);
      await loadData();
    } catch (err) {
      alert(err.message);
    } finally {
      setActionLoading(false);
    }
  };

  // Proses pengiriman atau pembatalan Sales Order
  const handleSalesOrderAction = async (id, action) => {
    const label = action === "ship" ? "mengirim" : "membatalkan";
    if (!window.confirm(`Yakin ingin ${label} Sales Order ini?`)) return;

    try {
      setActionLoading(true);
      const response = await fetch(`${API}/sales-orders/${id}/${action}`, {
        method: "POST",
      });
      const result = await response.json();
      if (!response.ok) {
        throw new Error(result.error || result.message || `Gagal ${label} Sales Order`);
      }

      alert(result.message || `Sales Order berhasil ${action === "ship" ? "dikirim" : "dibatalkan"}!`);
      await loadData();
    } catch (err) {
      alert(err.message);
    } finally {
      setActionLoading(false);
    }
  };

  // Menyelesaikan Work Order
  const handleCompleteWorkOrder = async (id) => {
    const confirmed = window.confirm(
      "Apakah kamu yakin ingin menyelesaikan Work Order ini?"
    );

    if (!confirmed) return;

    try {
      setActionLoading(true);

      const response = await fetch(`${API}/work-orders/${id}/complete`, {
        method: "POST",
      });

      const result = await response.json();

      if (!response.ok) {
        throw new Error(
          result.message || "Gagal menyelesaikan Work Order"
        );
      }

      alert("Work Order berhasil diselesaikan!");

      await loadData();
    } catch (err) {
      alert(err.message);
    } finally {
      setActionLoading(false);
    }
  };

  // Navigasi antarbagian dashboard
  const scrollToSection = (id) => {
    const section = document.getElementById(id);

    if (section) {
      const position =
        section.getBoundingClientRect().top + window.scrollY - 20;

      window.scrollTo({
        top: position,
        behavior: "smooth",
      });
    }
  };

  // Statistik dashboard
  const totalMaterials = materials.length;
  const totalProducts = products.length;
  const totalOrders = workOrders.length;
  const totalSalesOrders = salesOrders.length;

  const reservedProducts = products.reduce(
    (total, item) => total + Number(item.reserved || 0),
    0
  );

  const reservedMaterials = materials.reduce(
    (total, item) => total + Number(item.reserved || 0),
    0
  );

  return (
    <div className="app">
      {/* SIDEBAR */}
      <aside className="sidebar">
        <h2>
          Production<span>Hub</span>
        </h2>

        <nav>
          <a
            href="#dashboard"
            onClick={(e) => {
              e.preventDefault();
              scrollToSection("dashboard");
            }}
          >
            Dashboard
          </a>

          <a
            href="#inventory"
            onClick={(e) => {
              e.preventDefault();
              scrollToSection("inventory");
            }}
          >
            Inventory
          </a>

          <a
            href="#products"
            onClick={(e) => {
              e.preventDefault();
              scrollToSection("products");
            }}
          >
            Products
          </a>

          <a
            href="#sales-orders"
            onClick={(e) => {
              e.preventDefault();
              scrollToSection("sales-orders");
            }}
          >
            Sales Orders
          </a>

          <a
            href="#work-orders"
            onClick={(e) => {
              e.preventDefault();
              scrollToSection("work-orders");
            }}
          >
            Work Orders
          </a>
        </nav>

        <div className="sidebar-footer">
          Production Management System
        </div>
      </aside>

      {/* MAIN CONTENT */}
      <main className="main" id="dashboard">
        {/* HEADER */}
        <header className="header">
          <div>
            <h1>Dashboard</h1>
            <p>Production & Inventory Management</p>
          </div>

          <button onClick={loadData} disabled={loading}>
            Refresh Data ↻
          </button>
        </header>

        {error && <div className="error">{error}</div>}

        {loading ? (
          <p>Memuat data...</p>
        ) : (
          <>
            {/* STATISTICS */}
            <section className="stats">
              <div className="stat-card">
                <p>Total Materials</p>
                <h2>{totalMaterials}</h2>
              </div>

              <div className="stat-card">
                <p>Total Products</p>
                <h2>{totalProducts}</h2>
              </div>

              <div className="stat-card">
                <p>Total Work Orders</p>
                <h2>{totalOrders}</h2>
              </div>

              <div className="stat-card">
                <p>Reserved Materials</p>
                <h2>{reservedMaterials}</h2>
              </div>

              <div className="stat-card">
                <p>Reserved Products</p>
                <h2>{reservedProducts}</h2>
              </div>

              <div className="stat-card">
                <p>Total Sales Orders</p>
                <h2>{totalSalesOrders}</h2>
              </div>
            </section>

            {/* CREATE WORK ORDER */}
            <section className="panel">
              <h2>Buat Work Order</h2>

              <div className="form-group">
                <label>Pilih Produk</label>

                <select
                  value={selectedProduct}
                  onChange={(e) => {
                    setSelectedProduct(e.target.value);
                    setBomPreview(null);
                  }}
                >
                  <option value="">-- Pilih Produk --</option>

                  {products.map((product) => (
                    <option
                      key={product.product_id}
                      value={product.product_id}
                    >
                      {product.product_name}
                    </option>
                  ))}
                </select>
              </div>

              <div className="form-group">
                <label>Jumlah Produksi</label>

                <input
                  type="number"
                  min="1"
                  value={productionQuantity}
                  onChange={(e) => {
                    setProductionQuantity(e.target.value);
                    setBomPreview(null);
                  }}
                />
              </div>

              <button
                className="btn"
                onClick={handleBOMPreview}
                disabled={actionLoading}
              >
                {actionLoading
                  ? "Memproses..."
                  : "Preview Kebutuhan Bahan"}
              </button>

              {/* BOM PREVIEW */}
              {bomPreview && (
                <div style={{ marginTop: "20px" }}>
                  <h3>Hasil BOM Explosion</h3>

                  <p>
                    Produk: {bomPreview.product_name || "-"}
                  </p>

                  <p>
                    Versi BOM: {bomPreview.bom_version || "-"}
                  </p>

                  <p>
                    Jumlah Produksi:
                    {" "}
                    {bomPreview.production_quantity ||
                      productionQuantity}
                  </p>

                  <table>
                    <thead>
                      <tr>
                        <th>Bahan Baku</th>
                        <th>Kebutuhan</th>
                        <th>Satuan</th>
                      </tr>
                    </thead>

                    <tbody>
                      {(bomPreview.materials || []).map(
                        (material, index) => (
                          <tr key={index}>
                            <td>{material.material_name}</td>
                            <td>{material.total_quantity}</td>
                            <td>{material.unit}</td>
                          </tr>
                        )
                      )}
                    </tbody>
                  </table>

                  <button
                    className="btn"
                    style={{ marginTop: "15px" }}
                    onClick={handleCreateWorkOrder}
                    disabled={actionLoading}
                  >
                    {actionLoading
                      ? "Memproses..."
                      : "Buat Work Order"}
                  </button>
                </div>
              )}
            </section>

            {/* CREATE SALES ORDER */}
            <section className="panel" id="sales-orders">
              <h2>Buat Sales Order</h2>
              <div className="form-group">
                <label>Nama Pelanggan</label>
                <input
                  type="text"
                  value={salesCustomer}
                  onChange={(e) => setSalesCustomer(e.target.value)}
                  placeholder="Contoh: Pelanggan Demo"
                />
              </div>

              <div className="form-group">
                <label>Pilih Produk</label>
                <select value={salesProduct} onChange={(e) => setSalesProduct(e.target.value)}>
                  <option value="">-- Pilih Produk --</option>
                  {products.map((product) => (
                    <option key={product.product_id} value={product.product_id}>
                      {product.product_name} (Tersedia: {product.available})
                    </option>
                  ))}
                </select>
              </div>

              <div className="form-group">
                <label>Jumlah Pesanan</label>
                <input
                  type="number"
                  min="1"
                  value={salesQuantity}
                  onChange={(e) => setSalesQuantity(e.target.value)}
                />
              </div>

              <button className="btn" onClick={handleCreateSalesOrder} disabled={actionLoading}>
                {actionLoading ? "Memproses..." : "Buat Sales Order & Reservasi Stok"}
              </button>
            </section>

            {/* MATERIAL INVENTORY */}
            <section className="panel" id="inventory">
              <h2>Material Inventory</h2>

              <table>
                <thead>
                  <tr>
                    <th>Material</th>
                    <th>On Hand</th>
                    <th>Reserved</th>
                    <th>Available</th>
                  </tr>
                </thead>

                <tbody>
                  {materials.map((item) => (
                    <tr key={item.material_id}>
                      <td>
                        <strong>{item.material_name}</strong>
                        <small>{item.material_code}</small>
                      </td>

                      <td>
                        {item.on_hand} {item.base_unit}
                      </td>

                      <td>
                        {item.reserved} {item.base_unit}
                      </td>

                      <td className="available">
                        {item.available} {item.base_unit}
                      </td>
                    </tr>
                  ))}

                  {materials.length === 0 && (
                    <tr>
                      <td colSpan="4">Belum ada material.</td>
                    </tr>
                  )}
                </tbody>
              </table>
            </section>

            {/* PRODUCT INVENTORY */}
            <section className="panel" id="products">
              <h2>Finished Product Inventory</h2>

              <table>
                <thead>
                  <tr>
                    <th>Product</th>
                    <th>On Hand</th>
                    <th>Reserved</th>
                    <th>Available</th>
                  </tr>
                </thead>

                <tbody>
                  {products.map((item) => (
                    <tr key={item.product_id}>
                      <td>
                        <strong>{item.product_name}</strong>
                        <small>{item.product_code}</small>
                      </td>

                      <td>{item.on_hand}</td>

                      <td>{item.reserved}</td>

                      <td className="available">
                        {item.available}
                      </td>
                    </tr>
                  ))}

                  {products.length === 0 && (
                    <tr>
                      <td colSpan="4">Belum ada produk.</td>
                    </tr>
                  )}
                </tbody>
              </table>
            </section>

            {/* WORK ORDERS */}
<section className="panel" id="work-orders">
  <h2>Recent Work Orders</h2>

  <table>
    <thead>
      <tr>
        <th>Order Number</th>
        <th>Product</th>
        <th>Quantity</th>
        <th>BOM Version</th>
        <th>Status</th>
        <th>Action</th>
      </tr>
    </thead>

    <tbody>
      {workOrders.map((order) => (
        <tr key={order.id}>
          <td>{order.order_number}</td>

          <td>{order.product?.name || "-"}</td>

          <td>{order.quantity}</td>

          <td>V{order.bom_version}</td>

          <td>
            <span
              className={`badge ${String(
                order.status
              ).toLowerCase()}`}
            >
              {order.status}
            </span>
          </td>

          <td>
            {order.status === "RESERVED" ? (
              <button
                className="btn"
                onClick={() =>
                  handleCompleteWorkOrder(order.id)
                }
                disabled={actionLoading}
              >
                Complete
              </button>
            ) : (
              "-"
            )}
          </td>
        </tr>
      ))}

      {workOrders.length === 0 && (
        <tr>
          <td colSpan="6">Belum ada Work Order.</td>
        </tr>
      )}
    </tbody>
  </table>
</section>

            {/* SALES ORDERS */}
            <section className="panel">
              <h2>Sales Order</h2>
              <table>
                <thead>
                  <tr>
                    <th>Order Number</th>
                    <th>Pelanggan</th>
                    <th>Produk</th>
                    <th>Quantity</th>
                    <th>Status</th>
                    <th>Action</th>
                  </tr>
                </thead>
                <tbody>
                  {salesOrders.map((order) => (
                    <tr key={order.id}>
                      <td>{order.order_number}</td>
                      <td>{order.customer_name}</td>
                      <td>
                        {(order.items || []).map((item) => item.product?.name || `Product #${item.product_id}`).join(", ")}
                      </td>
                      <td>
                        {(order.items || []).reduce((sum, item) => sum + Number(item.quantity || 0), 0)}
                      </td>
                      <td>
                        <span className={`badge ${String(order.status).toLowerCase()}`}>
                          {order.status}
                        </span>
                      </td>
                      <td>
                        {order.status === "RESERVED" ? (
                          <>
                            <button
                              className="btn"
                              onClick={() => handleSalesOrderAction(order.id, "ship")}
                              disabled={actionLoading}
                            >
                              Ship
                            </button>{" "}
                            <button
                              className="btn"
                              onClick={() => handleSalesOrderAction(order.id, "cancel")}
                              disabled={actionLoading}
                            >
                              Cancel
                            </button>
                          </>
                        ) : "-"}
                      </td>
                    </tr>
                  ))}
                  {salesOrders.length === 0 && (
                    <tr><td colSpan="6">Belum ada Sales Order.</td></tr>
                  )}
                </tbody>
              </table>
            </section>
          </>
        )}
      </main>
    </div>
  );
}

export default App;