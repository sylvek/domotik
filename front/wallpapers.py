import urllib.request, json, re, os, pickle
from PIL import Image
from resizeimage import resizeimage
from tqdm import tqdm


# Variables
subreddits = ["earthporn"]  # which subreddits should we fetch?
timeranges = ["all"]        # top images from: all, year, month, week, day?
limit = 5000                # how many posts should we fetch?

output_list = "domotik/public/javascripts/data.js"
output_path = "domotik/public/img/"

min_width  = 1920           # minimal width, height and ratio to keep the image
min_height = 1080
min_ratio  =  1
resize     = True


# Create output directory if it doesn't exist
if not os.path.exists(output_path):
    os.makedirs(output_path)

# Load existing images list 
try:
    with open(output_list, 'r') as f:
        imglist = f.read().replace('var images = ', '')
        imglist = json.loads(imglist)
except:
    imglist = {}

# Load existing blacklist
try:
    with open ('blacklist', 'rb') as f:
        blacklist = pickle.load(f)
except:
    blacklist = []

last = ""

# Loop over subreddits and timeranges
for subreddit in subreddits:
    for timerange in timeranges:
        for i in tqdm(range(limit)):

            # Reddit API to fetch posts
            url = "https://www.reddit.com/r/" + subreddit + "/best/.json?t=" + timerange + "&limit=1" + "&after=" + last
            req = urllib.request.Request(url, headers={'User-Agent': 'Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:15.0) Gecko/20100101 Firefox/15.0.1'})

            with urllib.request.urlopen(req) as url:
                data = json.loads(url.read().decode())['data']['children']

                for post in data:
                    # Name, url
                    name = post['data']['name']
                    url = post['data']['url']
                    last = name

                    # Pass if already exists or in blacklist
                    if name in imglist or name in blacklist:
                        continue

                    # Catch Imgur URL
                    url = re.sub(r"^http(s)?:\/\/imgur\.com\/(.+)$", r"https://i.imgur.com/\2.jpg", url)

                    # Filter brackets, resolution and author from title
                    title = post['data']['title']
                    title = re.sub(r'\[.*\](.)*', '', title)
                    title = re.sub(r'\{.*\}(.)*', '', title)
                    title = re.sub(r'\(.*\)(.)*', '', title)
                    title = re.sub(r'(\d)+ ?x ?(\d)+(.)*', '', title)
                    title = re.sub(r'(by )*@\w+', '', title)
                    title = re.sub(r'(,|\.) +$', '', title)
                    tqdm.write("\n" + title)
                    tqdm.write(url)
                    
                    # Download image
                    try:
                        urllib.request.urlretrieve(url, output_path + name + ".jpg")
                
                    except:
                        blacklist.append(name) # we won't try to download it again
                        continue

                    # Check image resolution and ratio
                    test = False

                    try:
                        width, height = Image.open(output_path + name + ".jpg").size
                        test = (width >= min_width and height >= min_height and width/height > min_ratio)
                        tqdm.write(str(width) + "x" + str(height) + ", ratio " + str(width/height) + " => " + str(test))

                    except (KeyboardInterrupt, SystemExit):
                        raise

                    except:
                        test = False

                    if (test == False):
                        blacklist.append(name) # we won't try to download it again
                        os.remove(output_path + name + ".jpg")

                    else:
                        # Resize image
                        if resize == True:
                            with Image.open(output_path + name + ".jpg") as image:
                                img = resizeimage.resize_cover(image, [min_width, min_height])
                                img.save(output_path + name + ".jpg", image.format, quality=100)

                        imglist[name] = { 'title': title, 'file': name + ".jpg" }

    # Export image list
    with open(output_list, 'w') as f:
        print('var images = ' + json.dumps(imglist), file=f)

    # Export blacklist
    with open('blacklist', 'wb') as f:
        pickle.dump(blacklist, f)
