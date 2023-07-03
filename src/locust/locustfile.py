import time
#from locust import HttpUser, task, between
from locust import task, FastHttpUser, between
import random

#curl -X POST -H "Content-Type: application/json" --data 
#'{"team1" : "Korea","team2":"Italy","score":"0-0","phase":16}' http://$host:3000/match


class MatchesTrafficSimulation(FastHttpUser):
    wait_time = between(2, 5)
    countries = []
    uniques = set()
    l = 0
    @task
    def registerMatchPrediction(self):
        team1 = self.countries[random.randint(0,self.l)].replace("\n","")
        team2 = self.countries[random.randint(0,self.l)].replace("\n","")
        score = str(random.randint(0,5))+"-"+str(random.randint(0,5))
        phase = 2 ** random.randint(1,4)
        #print({"team1":team1,"team2":team2,"score":score,"phase":phase})
        if not ((team1+"-"+team2) in self.uniques):
        #    self.uniques.add(team1+"-"+team2)
            self.uniques.add(team1+"-"+team1)
            self.uniques.add(team2+"-"+team2)
            self.client.post("/match", json={"team1":team1,"team2":team2,"score":score,"phase":phase})

    def on_start(self):
        with open('countries.txt') as f:
            lines = f.readlines()

        for line in lines:
            self.countries.append(line)
        self.l = len(self.countries) - 1
