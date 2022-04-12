import sys
import time
def f(event):
    time.sleep(10)
    print("hellp")

if __name__ == "__main__":
    event = sys.argv[1]
    f(event)
