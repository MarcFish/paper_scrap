from selenium import webdriver
from selenium.webdriver.support.wait import WebDriverWait
from selenium.webdriver.support import expected_conditions as EC
from selenium.webdriver.common.by import By
from bs4 import BeautifulSoup as bs
import requests
import argparse
# from multiprocessing import Process, Queue
import json
import csv
import logging

parser = argparse.ArgumentParser()
parser.add_argument("--mode", type=str,default="journals", help="journals or conf")
parser.add_argument("--target", type=str, default="tkde", help="tkde,icde")
parser.add_argument("--volume", type=str, default="32", help="year or volume")

arg = parser.parse_args()

logging.basicConfig(filename="log.log", level=logging.INFO)

dblp_link = "https://dblp.org/db/" + arg.mode + "/" + arg.target + "/" + arg.target + arg.volume + ".html"
headers = {'User-Agent':"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/48.0.2564.116 Safari/537.36"}
r = requests.get(dblp_link, headers=headers).content
soup = bs(r, "lxml")
link_list = list()
for i in soup.select('.entry.article'):
    link_list.append(i.select('li.ee a')[0]['href'])

option = webdriver.ChromeOptions()
option.add_argument('headless')

logging.info(f"link:{dblp_link}")
logging.info(f"link count:{len(link_list)}")

def scrap_details(link_list):
    driver = webdriver.Chrome(options=option)
    driver.maximize_window()
    outs = list()
    for l in link_list:
        driver.get(l)
        try:
            locator = (By.TAG_NAME, "h1")
            WebDriverWait(driver, 30, 0.5).until(EC.presence_of_element_located(locator))
            title = driver.find_element_by_tag_name("h1").find_element_by_tag_name("span").text
            locator = (By.CLASS_NAME, "authors-info")
            WebDriverWait(driver, 30, 0.5).until(EC.presence_of_element_located(locator))
            authors = driver.find_elements_by_class_name("authors-info")
            author_list = list()
            for author in authors:
                author_list.append({
                    'name':author.find_element_by_tag_name("a").find_element_by_tag_name("span").text,
                    'id':author.find_element_by_tag_name("a").get_attribute("href"),
                })
            locator = (By.CLASS_NAME, "abstract-text")
            WebDriverWait(driver, 30, 0.5).until(EC.presence_of_element_located(locator))
            abstract = driver.find_element_by_class_name("abstract-text").find_element_by_tag_name("div").find_element_by_tag_name("div").find_element_by_tag_name("div").text

            locator = (By.ID, "references-header")
            WebDriverWait(driver, 30, 0.5).until(EC.presence_of_element_located(locator))
            references = driver.find_element_by_id("references-header")
            driver.execute_script("arguments[0].click();", references)
            references = driver.find_element_by_id("references-section-container").find_elements_by_class_name("reference-container")
            ref_list = list()
            for ref in references:
                ref_list.append(ref.find_element_by_class_name("row").find_element_by_class_name("col-12").find_elements_by_tag_name("span")[1].get_attribute("innerHTML"))
        except:
            logging.error(f"wrong link:{l}")
            continue
        out = [l, title, abstract, json.dumps(author_list), json.dumps(ref_list)]
        outs.append(out)
    driver.close()
    return outs


def save(fname, outs):
    header = ['doi','title', 'abstract', 'authors','references']
    with open(fname, 'w', encoding='utf-8', newline='') as f:
        writer = csv.writer(f)
        writer.writerow(header)
        for out in outs:
            writer.writerow(out)


if __name__ == '__main__':
    outs = scrap_details(link_list)
    save("."+arg.mode+"_"+arg.target+"_"+arg.volume+".csv", outs)
