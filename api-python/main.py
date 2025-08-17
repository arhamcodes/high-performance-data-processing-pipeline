from fastapi import FastAPI, HTTPException
from pydantic import BaseModel
from typing import List, Optional, Dict
import asyncpg
import aio_pika
import json
from datetime import datetime
import uuid

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

# @app.post("/ingest/")
# async def ingest_order(order: Order):
#     print("Received order data:", order)
#     return {"status": "success", "order": order}

# @app.get("/health")
# async def health_check():
#     return {"status": "healthy"}


class OrderStatus(BaseModel):
    id: str
    timestamp: str
    status: str
    order_data: Dict

# @app.post("/ingest/")
# async def ingest_order(order: Order):
#     # Generate unique ID for tracking
#     # order_id = str(uuid.uuid4())
#     # timestamp = datetime.utcnow().isoformat()
#     order_id = str(uuid.uuid4())
#     timestamp = datetime.utcnow()
    
#     # Create database record
#     try:
#         conn = await asyncpg.connect('postgresql://arham@postgres:5432/orders')
#         await conn.execute('''
#             INSERT INTO orders (id, timestamp, status, order_data)
#             VALUES ($1, $2, $3, $4)
#         ''', order_id, timestamp, 'pending', order.dict())
#         await conn.close()
#     except Exception as e:
#         raise HTTPException(status_code=500, detail=f"Database error: {str(e)}")


@app.post("/ingest/")
async def ingest_order(order: Order):
    order_id = str(uuid.uuid4())
    timestamp = datetime.utcnow()

    try:
        conn = await asyncpg.connect('postgresql://arham@postgres:5432/orders')
        await conn.execute('''
            INSERT INTO orders (id, timestamp, status, order_data)
            VALUES ($1, $2, $3, $4)
        ''', order_id, timestamp, 'pending', json.dumps(order.dict()))
        await conn.close()
    except Exception as e:
        raise HTTPException(status_code=500, detail=f"Database error: {str(e)}")
    
    # Publish to RabbitMQ
    # try:
    #     connection = await aio_pika.connect_robust("amqp://guest:guest@rabbitmq/")
    #     channel = await connection.channel()
    #     queue = await channel.declare_queue("orders")
        
    #     message = {
    #         "id": order_id,
    #         "timestamp": timestamp,
    #         "order": order.dict()
    #     }
        
    #     await channel.default_exchange.publish(
    #         aio_pika.Message(body=json.dumps(message).encode()),
    #         routing_key="orders"
    #     )
    #     await connection.close()
    # except Exception as e:
    #     raise HTTPException(status_code=500, detail=f"Queue error: {str(e)}")
    try:
        connection = await aio_pika.connect_robust("amqp://guest:guest@rabbitmq/")
        channel = await connection.channel()
        queue = await channel.declare_queue("orders")

        message = {
            "id": order_id,
            "timestamp": timestamp.isoformat(),
            "order": order.dict()
        }

        await channel.default_exchange.publish(
            aio_pika.Message(body=json.dumps(message).encode()),
            routing_key="orders"
        )
        await connection.close()
    except Exception as e:
        raise HTTPException(status_code=500, detail=f"Queue error: {str(e)}")

    return {"status": "success", "order_id": order_id}

@app.get("/status/{order_id}")
async def get_order_status(order_id: str):
    try:
        conn = await asyncpg.connect('postgresql://arham@postgres:5432/orders')
        record = await conn.fetchrow(
            'SELECT id, timestamp, status, order_data FROM orders WHERE id = $1',
            order_id
        )
        await conn.close()
        
        if not record:
            raise HTTPException(status_code=404, detail="Order not found")
            
        return OrderStatus(
            id=record['id'],
            timestamp=record['timestamp'],
            status=record['status'],
            order_data=record['order_data']
        )
    except Exception as e:
        raise HTTPException(status_code=500, detail=f"Database error: {str(e)}")