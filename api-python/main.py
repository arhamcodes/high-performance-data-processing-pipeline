from fastapi import FastAPI
from pydantic import BaseModel
from typing import List, Optional
app = FastAPI()

class Address(BaseModel):
    street: str
    city: str
    state: str
    zip_code: str
    country: str

class Customer(BaseModel):
    id: str
    email: str
    firstName: str
    lastName: str
    shippingAddress: Address  
    billingAddress: Address   

class Item(BaseModel):
    productId: str
    name: str
    price: float
    quantity: int
    variant: Optional[str] = None 

class ShippingMethod(BaseModel):
    id: str
    name: str
    cost: float

class Payment(BaseModel):
    method: str
    token: str
    amount: float
    currency: str

class Order(BaseModel):
    customer: Customer
    items: List[Item]  
    shippingMethod: ShippingMethod
    payment: Payment
    orderTotal: float
    taxAmount: float
    discountAmount: float
    notes: Optional[str] = None 

@app.post("/ingest/")
async def ingest_order(order: Order):
    print("Received order data:", order)
    return {"status": "success", "order": order}

@app.get("/health")
async def health_check():
    return {"status": "healthy"}