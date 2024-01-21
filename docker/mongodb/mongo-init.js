db = db.getSiblingDB('cloud');
db.createCollection("run_time")

db.run_time.insert({
    started_at: new Date()
})