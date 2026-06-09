const collectionNames = db.getCollectionNames();

if (!collectionNames.includes("parts")) {
  db.createCollection("parts", {
    validator: {
      $jsonSchema: {
        bsonType: "object",
        required: [
          "schema_version",
          "name",
          "description",
          "price_cents",
          "stock_quantity",
          "category",
          "dimensions",
          "manufacturer",
          "tags",
          "metadata",
          "created_at",
          "updated_at",
        ],
        properties: {
          schema_version: {
            bsonType: "int",
            enum: [1],
          },
          name: {
            bsonType: "string",
          },
          description: {
            bsonType: "string",
          },
          price_cents: {
            bsonType: ["int", "long"],
            minimum: 0,
          },
          stock_quantity: {
            bsonType: ["int", "long"],
            minimum: 0,
          },
          category: {
            bsonType: "string",
            enum: ["ENGINE", "FUEL", "PORTHOLE", "WING"],
          },
          dimensions: {
            bsonType: "object",
            required: ["length", "width", "height", "weight"],
          },
          manufacturer: {
            bsonType: "object",
            required: ["name", "country", "website"],
          },
          tags: {
            bsonType: "array",
            items: {
              bsonType: "string",
            },
          },
          metadata: {
            bsonType: "object",
          },
          created_at: {
            bsonType: "date",
          },
          updated_at: {
            bsonType: "date",
          },
        },
      },
    },
  });
}

db.parts.createIndexes([
  { key: { name: 1 }, name: "parts_name_idx" },
  { key: { category: 1 }, name: "parts_category_idx" },
  {
    key: { "manufacturer.country": 1 },
    name: "parts_manufacturer_country_idx",
  },
  { key: { tags: 1 }, name: "parts_tags_idx" },
]);
